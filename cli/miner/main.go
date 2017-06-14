package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gorilla/websocket"

	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions = make([]*core.Transaction, 0)
	consensus           *c.Consensus
	network             *c.Network
	block               *core.Block      // pending block
	blockchain          *core.Blockchain // blockchain
	localID             string           // id of this peer, used to grant PoW
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	httpPort := flag.Int("p", 3000, "HTTP listen port")
	rsaFile := flag.String("i", "../private.pem", "RSA key file")
	blockchainFile := flag.String("b", "./blockchain", "Blockchain file")
	flag.Parse()

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFile)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}

	// hash private key to get id
	hash := sha256.New()
	io.WriteString(hash, fmt.Sprintf("%v", rsaPrivateKey))
	localID = fmt.Sprintf("%x", hash.Sum(nil))

	// Genesis Peer
	if flag.NArg() < 1 {
		consensus = c.NewConsensus()
		network = c.NewNetwork()
		blockchain = core.NewBlockchain()
	} else {
		// args consensus port
		p2pPort, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			log.Fatalf("given p2p port could not be parsed as int: %v: %v", p2pPort, err)
		}
		// Handshake protocol: Connect and get dial response
		addr := fmt.Sprintf("ws://localhost:%d/ws", p2pPort)
		log.Printf("Connecting to %s", addr)
		dialResponse, err := c.Connect(addr)
		if err != nil {
			log.Printf("disconnected from %s : %v", addr, err)
		}
		consensus = dialResponse.Consensus
		network = dialResponse.Network
		blockchain = new(core.Blockchain)
		err = core.Load(*blockchainFile, blockchain)
		if err != nil {
			blockchain = core.NewBlockchain()
			log.Printf("could not load blockchain stored at %s: %v", *blockchainFile, err)
		}
		blockchain = blockchain.Fetch(dialResponse.Blockchain)
		// Register this peer because it is not self registered in its network
		peer := &c.Peer{Conn: dialResponse.Conn}
		network.Add(peer)
	}

	// Serve HTTP
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/ws", handleConnection)
	go func() {
		log.Printf("HTTP listening to port %d", *httpPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
	}()

	go func() {
		// Listen and serve peers
		for _, peer := range network.Peers {
			go ListenAndServe(peer)
		}
	}()

	// mine
	go Mine()

	// Capture SIGTERM
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Printf("Saving blockchain locally at %s", *blockchainFile)
		core.Save(*blockchainFile, blockchain)
		cleanupDone <- true
	}()
	<-cleanupDone
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	var t c.TransactionMessage
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// verify
	err = crypto.Verify(t.Signature, t.Hash, t.RsaPublicKey)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Printf("%#v", t.Transaction)
	if block.IsTransactionValid(t.Transaction) {
		pendingTransactions = append(pendingTransactions, t.Transaction)
		msg, _ := json.Marshal(t)
		network.Broadcast(c.Message{Type: c.Transaction, Message: msg})
	}
}

func handleBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain)
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("could not upgrade websocket", err)
	}
	defer conn.Close()

	dialResponse := &c.DialResponse{Consensus: consensus, Network: network, Blockchain: blockchain}
	// send dial response
	conn.WriteJSON(dialResponse)
	// Register our new peer
	peer := &c.Peer{Conn: conn}
	network.Add(peer)
	// Serve this new peer
	go ListenAndServe(peer)
}

// Mine work on proof-of-work on the last block of its blockchain
func Mine() {
	block = blockchain.GetLastBlock()
	for {
		// Proof of Work
		for !crypto.MatchHash(block.ToHash(), consensus.Difficulty) {
			block = blockchain.Mine(consensus.Difficulty)
		}
		log.Println("BLOCK MINED !")
		// Present block
		log.Println("broadcasting to ", network)
		msg, _ := json.Marshal(&c.BlockMessage{Block: block, From: localID, Signature: false})
		network.Broadcast(c.Message{Type: c.Block, Message: msg})
		// aggregate peers responses
		responses := network.Aggregate()
		// check validity
		valid := consensus.Validate(block, responses)
		if valid {
			blockchain.AppendBlock(block)
			msg, _ := json.Marshal(&c.BlockPoWMessage{Block: block, From: localID, Signatures: responses})
			network.Broadcast(c.Message{Type: c.BlockPoW, Message: msg})
		}
		blockchain.AppendBlock(block)
		// update next tick
		consensus.UpdateNext()
		// work on "new clean" block: link + transactions history
		block = blockchain.GenNext(pendingTransactions)
		// flush pending transactions
		pendingTransactions = make([]*core.Transaction, 0)
	}
}

func ListenAndServe(peer *c.Peer) {
	log.Println("listening to ", peer)
	for {
		var msg c.Message
		// Read in a new message as JSON and map it to a Message object
		err := peer.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		// handling message
		switch msg.Type {
		case c.PeerStatus:
			log.Println("msg PeerStatus")
			var peerStatusMsg *c.PeerStatusMessage
			err = json.Unmarshal(msg.Message, &peerStatusMsg)
			network.Add(peerStatusMsg.Peer)
			go ListenAndServe(peerStatusMsg.Peer)
		case c.Transaction:
			log.Println("msg Transaction")
			var transactionMsg *c.TransactionMessage
			err = json.Unmarshal(msg.Message, &transactionMsg)
			if block.IsTransactionValid(transactionMsg.Transaction) {
				pendingTransactions = append(pendingTransactions, transactionMsg.Transaction)
			}
		case c.Block:
			log.Println("msg Block")
			var blockMsg *c.BlockMessage
			err = json.Unmarshal(msg.Message, &blockMsg)
			if blockchain.IsBlockValid(blockMsg.Block) {
				msg, _ := json.Marshal(&c.BlockMessage{Block: blockMsg.Block, From: blockMsg.From, Signature: true})
				peer.Conn.WriteJSON(c.Message{Type: c.Block, Message: msg})
			}
		case c.BlockPoW:
			log.Println("msg BlockPoW")
			var blockPoWMsg *c.BlockPoWMessage
			err = json.Unmarshal(msg.Message, &blockPoWMsg)
			blockchain.AppendBlock(block)
			// update next tick
			consensus.UpdateNext()
			// work on "new clean" block: link + transactions history
			block = blockchain.GenNext(pendingTransactions)
			// flush pending transactions
			pendingTransactions = make([]*core.Transaction, 0)
		}
	}
}
