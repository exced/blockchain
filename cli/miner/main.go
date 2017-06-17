package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"

	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions = make([]*core.Transaction, 0)
	mutex               = &sync.Mutex{}
	consensus           *c.Consensus
	network             *c.Network
	block               *core.Block      // pending block
	blockchain          *core.Blockchain // blockchain
	localID             string           // id of this peer, used to grant PoW and authenticate to the network
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
	blockchainFile := flag.String("b", "./blockchain.bc", "Blockchain file")
	flag.Parse()

	// miner ID
	localID, err := crypto.RsaID(*rsaFile)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}

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
		selfPeer := &c.Peer{Address: fmt.Sprintf("ws://localhost:%d/ws", *httpPort)}
		addr := fmt.Sprintf("ws://localhost:%d/ws", p2pPort)
		log.Printf("Connecting to %s", addr)
		conn, dialResponse, err := c.Connect(addr, selfPeer)
		if err != nil {
			log.Printf("disconnected from %s : %v", addr, err)
		}
		consensus = dialResponse.Consensus
		network = dialResponse.Network

		// load blockchain
		blockchain = new(core.Blockchain)
		err = core.Load(*blockchainFile, blockchain)
		if err != nil {
			blockchain = core.NewBlockchain()
			log.Printf("could not load blockchain stored at %s: %v", *blockchainFile, err)
		}
		blockchain = blockchain.Fetch(dialResponse.Blockchain)

		// Listen and serve peers
		for _, peer := range network.Peers {
			conn, _, err := c.Connect(peer.Address, selfPeer)
			if err != nil {
				log.Printf("disconnected from %s : %v", addr, err)
			}
			peer.Conn = conn
			go ListenAndServe(conn)
		}
		// Register this peer because it is not self registered in its network
		peer := &c.Peer{Address: addr, Conn: conn}
		network.Add(peer)
		go ListenAndServe(peer.Conn)
	}

	// current working block
	block = blockchain.GetLastBlock()

	// Serve HTTP
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/withdraw", handleWithdraw)
	http.HandleFunc("/deposit", handleDeposit)
	http.HandleFunc("/ws", handleConnection)
	go func() {
		log.Printf("HTTP listening to port %d", *httpPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
	}()

	// mine
	go Mine()

	// save blockchain if interrupt signal is captured
	core.SaveOnInterrupt(*blockchainFile, blockchain)
}

// handleWithdraw broadcast given transaction from client to network
func handleWithdraw(w http.ResponseWriter, r *http.Request) {
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
		log.Println("Broadcast transaction to ", network)
		network.Broadcast(c.Message{Type: c.Transaction, Message: msg})
	}
}

// handleDeposit broadcast given transaction from client to network
func handleDeposit(w http.ResponseWriter, r *http.Request) {
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
	pendingTransactions = append(pendingTransactions, t.Transaction)
	msg, _ := json.Marshal(t)
	log.Println("Broadcast transaction to ", network)
	network.Broadcast(c.Message{Type: c.Transaction, Message: msg})
}

// handleBlockchain retrieves a copy of current blockchain
func handleBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(blockchain)
}

// handleConnection add this new peer to network and listen and serve
func handleConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade websocket", err)
	}
	defer conn.Close()
	// read peer msg
	var msg c.Message
	// Read in a new message as JSON and map it to a Message object
	err = conn.ReadJSON(&msg)
	if err != nil {
		log.Printf("error read Message on conn %#v: %v", conn, err)
	}
	if msg.Type != c.PeerConn {
		log.Fatalf("wrong msg sent by peer %#v : %v", msg, err)
	}
	var peerConnMsg *c.PeerConnMessage
	err = json.Unmarshal(msg.Message, &peerConnMsg)
	peer := peerConnMsg.Peer
	peer.Conn = conn
	if err != nil {
		log.Println("could not read peer connection", err)
	}
	// Add new peer to network
	dialResponse := &c.DialResponse{Consensus: consensus, Network: network, Blockchain: blockchain}
	// send dial response
	conn.WriteJSON(dialResponse)
	// Register our new peer
	network.Add(peer)
	// Serve this new peer
	ListenAndServe(conn)
}

// Mine work on proof-of-work on the last block of its blockchain
func Mine() {
	for {
		// Proof of Work
		for !crypto.MatchHash(block.ToHash(), consensus.Difficulty) {
			mutex.Lock()
			block = blockchain.Mine(consensus.Difficulty)
			mutex.Unlock()
		}
		mutex.Lock()
		// Present block
		msg, _ := json.Marshal(&c.BlockMessage{Block: block, From: localID, Signature: false})
		network.Broadcast(c.Message{Type: c.Block, Message: msg})
		log.Println("BLOCK MINED !")
		// aggregate peers responses
		responses := network.Aggregate()
		log.Println("responses aggregate ", responses)
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
		mutex.Unlock()
	}
}

// ListenAndServe listen and serve messages from websocket connection
func ListenAndServe(conn *websocket.Conn) {
	for {
		var msg c.Message
		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error read Message on conn %#v: %v", conn, err)
		}
		// handling message
		switch msg.Type {
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
				conn.WriteJSON(c.Message{Type: c.Block, Message: msg})
			}
		case c.BlockPoW:
			log.Println("msg BlockPoW")
			var blockPoWMsg *c.BlockPoWMessage
			err = json.Unmarshal(msg.Message, &blockPoWMsg)
			mutex.Lock()
			blockchain.AppendBlock(block)
			// update next tick
			consensus.UpdateNext()
			// work on "new clean" block: link + transactions history
			block = blockchain.GenNext(pendingTransactions)
			// flush pending transactions
			pendingTransactions = make([]*core.Transaction, 0)
			mutex.Unlock()
		}
	}
}
