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
	blockchain          *core.Blockchain
	network             = c.NewNetwork()
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
	personalID := fmt.Sprintf("%x", hash.Sum(nil))

	log.Println(personalID)

	// Genesis Peer
	if flag.NArg() < 1 {
		// New Consensus
		consensus = c.NewConsensus()
		blockchain = core.NewBlockchain()
	} else {
		// Connect and Get Consensus
		p2pPort, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			log.Fatalf("given p2p port could not be parsed as int: %v: %v", p2pPort, err)
		}
		log.Print(fmt.Sprintf("Connecting to ws://localhost:%d/ws", p2pPort))
		consensus, err = c.Connect(fmt.Sprintf("ws://localhost:%d/ws", p2pPort))
		if err != nil {
			log.Printf("disconnected from %d : %v", p2pPort, err)
		}
		log.Printf("consensus: %#v ", consensus)
	}

	// Serve HTTP
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/send", handleSend)
	http.HandleFunc("/ws", handleConnection)
	go func() {
		log.Printf("HTTP listening to port %d", *httpPort)
		http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)
	}()

	// mine
	go mine()

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

	// Register our new peer
	network.Peers[conn] = "conn"

	for {
		var msg c.Message
		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(network.Peers, conn)
			break
		}
		// handling message
		switch msg.Type {
		case c.PeerStatus:
			log.Println("msg PeerStatus")
		case c.Transaction:
			log.Println("msg Transaction")
		case c.Block:
			log.Println("msg Block")
		}
		// Send the newly received message to the network
		network.Broadcast(msg)
	}
}

func mine() {
	var b *core.Block
	// notify
	go func() {
		for range time.Tick(time.Until(consensus.Tick)) {
		}
		consensus.Tick = consensus.Tick.Add(consensus.HashRate)
	}()
	for {

	}
	b = blockchain.Mine(consensus.Difficulty)
	if crypto.MatchHash(b.ToHash(), consensus.Difficulty) {

	}
	// Broadcast Mined Block
	network.Broadcast(c.Message{Type: c.Block, Message: &c.BlockMessage{Block: b}})
}
