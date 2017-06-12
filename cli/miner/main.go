package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	pb "github.com/exced/blockchain/cli/api"
	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*pb.TransactionMessage
	consensus           *c.Consensus
)

func main() {
	p2pPort := flag.Int("p2p", 6000, "P2P listen address")
	httpPort := flag.Int("http", 8000, "HTTP port")
	rsaFilePath := flag.String("r", "../private.pem", "RSA key file")
	genesisMode := flag.Bool("g", false, "Genesis mode")
	// blockchainFilePath := flag.String("b", "./blockchain", "Blockchain file")
	flag.Parse()

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey
	fmt.Println("rsapublickey ", rsaPublicKey)

	// Genesis Peer
	if *genesisMode {
		consensus = c.NewConsensus()
	} else {
		consensus, err = c.Connect(fmt.Sprintf("localhost:%d", *p2pPort))
		if err != nil {
			log.Fatalf("could not connect to consensus %d : %v", *p2pPort, err)
		}
	}

	// Handle peers connection
	// http.HandleFunc("/ws", websocket.Handler(consensus.HandlePeerConnection))
	http.HandleFunc("/blockchain", handleBlockchain)
	log.Printf("HTTP listening to port %d", *httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil)

	// Serve Network
	// go consensus.ListenAndServe()

	// mine
	go func() {
		// b := c.Blockchain.Mine(c.Difficulty)
		// if crypto.MatchHash(b.ToHash(), c.Difficulty) {

		// }
		// Broadcast Mined Block
		// broadcast <- &Message{Type: BlockMessage, Message: &BlockMessage{Block: b}}
	}()
}

func handleBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(consensus.Blockchain)
}
