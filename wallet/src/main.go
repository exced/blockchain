package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/exced/blockchain/core"
	"github.com/gorilla/websocket"
)

var (
	blockchain   *core.Blockchain     // blockchain
	transactions []*core.Transaction  // pending transactions
	consensus    *consensus.Consensus // consensus
	upgrader     websocket.Upgrader
)

func main() {
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	p2pAddr := flag.String("p2p", ":6000", "P2P listen address")
	flag.Parse()

	// fetch blockchain
	// take local saved blockchain
	blockchain = core.NewBlockchain()
	blockchain = consensus.Fetch(blockchain)

	// HTTP API
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Client Logic
	http.HandleFunc("/deposit", handleDeposit)
	http.HandleFunc("/withdraw", handleWithdraw)

	// http serve
	log.Println("http server started on", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal("Could not serve http: ", err)
	}
}

// handleConnections broadcast newly connected or disconnected peer to peers
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register our new peer
	peers[ws] = true

	msg := &PeerConnection{
		Conn:   ws,
		Status: true,
	}

	// Send it out to every peer that is currently connected
	for peer := range peers {
		err := peer.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			peer.Close()
			delete(peers, peer)
		}
	}
}

func handleDeposit(w http.ResponseWriter, r *http.Request) {
	// decode transaction
	decoder := json.NewDecoder(r.Body)
	var t core.Transaction
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	// add transaction to pending
	transactions = append(transactions, &t)
}

func handleWithdraw(w http.ResponseWriter, r *http.Request) {
	// decode transaction
	decoder := json.NewDecoder(r.Body)
	var t core.Transaction
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}

	// add transaction to pending
	transactions = append(transactions, &t)
}
