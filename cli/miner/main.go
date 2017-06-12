package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*core.Transaction
	consensus           *c.Consensus
)

func main() {
	httpPort := flag.Int("http", 3000, "HTTP port")
	rsaFilePath := flag.String("i", "../private.pem", "RSA key file")
	// blockchainFilePath := flag.String("b", "./blockchain", "Blockchain file")
	flag.Parse()

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey
	log.Println("rsapublickey ", rsaPublicKey)

	// Genesis Peer
	if flag.NArg() < 1 {
		consensus = c.NewConsensus()
	} else {
		// args
		p2pPort, err := strconv.ParseInt(flag.Arg(0), 10, 64)
		if err != nil {
			log.Fatalf("given p2p port could not be parsed as int: %v: %v", p2pPort, err)
		}
		consensus, err = c.Connect(fmt.Sprintf("localhost:%d", p2pPort))
		if err != nil {
			log.Fatalf("could not connect to consensus %d : %v", p2pPort, err)
		}
	}

	// Handle peers connection
	// http.HandleFunc("/ws", websocket.Handler(consensus.HandlePeerConnection))
	http.HandleFunc("/blockchain", handleBlockchain)
	http.HandleFunc("/send", handleSend)

	// Serve HTTP
	log.Printf("HTTP listening to port %d", *httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)

	// Serve Network
	go consensus.ListenAndServe()

	// mine
	go func() {
		// b := c.Blockchain.Mine(c.Difficulty)
		// if crypto.MatchHash(b.ToHash(), c.Difficulty) {

		// }
		// Broadcast Mined Block
		// broadcast <- &Message{Type: BlockMessage, Message: &BlockMessage{Block: b}}
	}()
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
	json.NewEncoder(w).Encode(consensus.Blockchain)
}
