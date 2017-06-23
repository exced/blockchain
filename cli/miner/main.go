package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

var (
	consensus    *c.Consensus       // consensus object
	network      *c.Network         // network object
	transactions *core.Transactions // pending transactions
	blockchain   *core.Blockchain   // blockchain
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

	// Genesis Peer
	if flag.NArg() < 1 {
		consensus = c.NewConsensus()
		network = c.NewNetwork()
		blockchain = core.NewBlockchain()
	} else {
		// args consensus port
		p2pPort, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			log.Fatalf("p2p port could not be parsed as int: %v: %v", p2pPort, err)
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

	// current working transactions
	transactions = blockchain.GetLastBlock().Transactions

	// Serve HTTP
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/transaction", handleTransaction)
	http.HandleFunc("/ws", handleConnection)
	go func() {
		log.Printf("HTTP listening to port %d", *httpPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
	}()

	// miner key
	localKey, err := crypto.RsaID(*rsaFile)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	// mine for local key
	go Mine(localKey)

	// save blockchain if interrupt signal is captured
	core.SaveOnInterrupt(*blockchainFile, blockchain)
}

// handleTransaction broadcast given transaction from client to network
func handleTransaction(w http.ResponseWriter, r *http.Request) {
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
	log.Printf("request: %#v", t.Transaction)
	json.NewEncoder(w).Encode(blockchain.GetLastBlock().Transactions.IsValid(t.Transaction))
	transactions.Append(t.Transaction)
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
func Mine(localKey string) {
	for {
		// Proof of Work
		for !blockchain.PoW(consensus.Difficulty) {
			blockchain.Mine()
		}
		// Present block
		msg, _ := json.Marshal(&c.BlockMessage{Block: blockchain.GetLastBlock(), From: localKey, Signature: false})
		network.Broadcast(c.Message{Type: c.Block, Message: msg})
		// aggregate peers responses
		responses := network.Aggregate()
		// check validity
		if network.Validate(blockchain.GetLastBlock(), responses) {
			blockchain.AppendBlock(blockchain.GetLastBlock())
			msg, _ = json.Marshal(&c.BlockPoWMessage{Block: blockchain.GetLastBlock(), From: localKey, Signatures: responses})
			network.Broadcast(c.Message{Type: c.BlockPoW, Message: msg})
		}
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
			transactions.Append(transactionMsg.Transaction)
			log.Println("transactions: ", transactions)
		case c.Block:
			log.Println("msg Block")
			var blockMsg *c.BlockMessage
			err = json.Unmarshal(msg.Message, &blockMsg)
			var answer *c.BlockMessage
			if blockchain.IsBlockValid(blockMsg.Block) {
				answer = &c.BlockMessage{Block: blockMsg.Block, From: blockMsg.From, Signature: true}
			} else {
				answer = &c.BlockMessage{Block: blockMsg.Block, From: blockMsg.From, Signature: false}
			}
			msg, _ := json.Marshal(answer)
			conn.WriteJSON(c.Message{Type: c.Block, Message: msg})
		case c.BlockPoW:
			log.Println("msg BlockPoW")
			var blockPoWMsg *c.BlockPoWMessage
			err = json.Unmarshal(msg.Message, &blockPoWMsg)
			// check validity
			if network.Validate(blockPoWMsg.Block, blockPoWMsg.Signatures) {
				blockchain.AppendBlock(blockPoWMsg.Block)
				// reward
				reward := core.NewTransaction(fmt.Sprintf("%v", blockPoWMsg.Signatures), blockPoWMsg.From, 1)
				transactions.Append(reward)
			}
		}
	}
}
