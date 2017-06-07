package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/exced/simple-blockchain/core"
	"github.com/gorilla/websocket"
)

var (
	blockchain *core.Blockchain         // blockchain
	peers      map[*websocket.Conn]bool // connected peers
	broadcast  chan PeerMessage         // broadcast channel
	upgrader   websocket.Upgrader
)

// PeerMessage defines our peer message object
type PeerMessage struct {
	Block string `json:"block"`
	State string `json:"state"`
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register our new peer
	peers[ws] = true

	for {
		var msg PeerMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(peers, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for peer := range peers {
			err := peer.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				peer.Close()
				delete(peers, peer)
			}
		}
	}
}

func main() {
	blockchain = core.NewBlockchain()
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	p2pAddr := flag.String("p2p", ":6000", "p2p server address.")
	flag.Parse()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming peer messages
	go handleMessages()

	// http serve
	log.Println("http server started on", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal("Could not serve http: ", err)
	}
}
