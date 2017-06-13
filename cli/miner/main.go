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
	"time"

	"github.com/gorilla/websocket"

	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*core.Transaction
	consensus           *c.Consensus
	network             *c.Network
	blockchain          *core.Blockchain
	localID             string // id of the local peer
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
		// Handshake protocol: Connect and Get Consensus and currnet network
		log.Print(fmt.Sprintf("Connecting to ws://localhost:%d/ws", p2pPort))
		consensus, network, err = c.Connect(fmt.Sprintf("ws://localhost:%d/ws", p2pPort))
		if err != nil {
			log.Printf("disconnected from %d : %v", p2pPort, err)
		}
	}

	// Serve HTTP
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/ws", handleConnection)
	go func() {
		log.Printf("HTTP listening to port %d", *httpPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
	}()

	// Listen and serve peers
	for _, peer := range network.Peers {
		go ListenAndServe(peer)
	}

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
	pendingTransactions = append(pendingTransactions, t.Transaction)
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
	// send my consensus
	conn.WriteJSON(consensus)
	// send my network
	conn.WriteJSON(network)
	// Register our new peer
	peer := &c.Peer{Conn: conn}
	network.Add(peer)
	// Serve this new peer
	go ListenAndServe(peer)
}

func Mine() {
	// sleep until next Tick to begin mining
	time.Sleep(time.Until(consensus.Next))
	var b *core.Block
	// Mine block
	go func() {
		for !crypto.MatchHash(b.ToHash(), consensus.Difficulty) {
			b = blockchain.Mine(consensus.Difficulty)
		}
	}()
	// Broadcast mined block
	go func() {
		for range time.Tick(consensus.HashRate) {
			if crypto.MatchHash(b.ToHash(), consensus.Difficulty) {
				log.Println("broadcasting to ", network)
				msg, _ := json.Marshal(&c.BlockMessage{Block: b, From: localID})
				network.Broadcast(c.Message{Type: c.Block, Message: msg})
			}
			// append accepted block
			accepted := consensus.Accepted()
			blockchain.AppendBlock(accepted)
			// update next tick
			consensus.UpdateNext()
		}
	}()
}

func ListenAndServe(peer *c.Peer) {
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
		case c.Transaction:
			log.Println("msg Transaction")
			var transaction *core.Transaction
			err = json.Unmarshal(msg.Message, &transaction)
			pendingTransactions = append(pendingTransactions, transaction)
		case c.Block:
			log.Println("msg Block")
			var block *core.Block
			err = json.Unmarshal(msg.Message, &block)
		}
	}
}
