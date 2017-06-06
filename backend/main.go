package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan PeerMessage)       // broadcast channel

// PeerMessage defines our peer message object
type PeerMessage struct {
	Block string `json:"block"`
	Stage string `json:"stage"`
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg PeerMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func main() {
	httpAddr := flag.String("api", ":3001", "api server address.")
	p2pAddr := flag.String("p2p", ":6001", "p2p server address.")
	initialPeers := flag.String("peers", "ws://localhost:6001", "initial peers")
	flag.Parse()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)
}
