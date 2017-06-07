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
	State string `json:"state"`
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

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// handleDeposit add deposit transaction to current pendings transactions
func handleDeposit(w http.ResponseWriter, r *http.Request) {

}

// handleWithdraw add withdraw transaction to current pendings transactions
func handleWithdraw(w http.ResponseWriter, r *http.Request) {

}

func main() {
	httpAddr := flag.String("http", ":3001", "HTTP listen address")
	// p2pAddr := flag.String("p2p", ":6001", "p2p server address.")
	// initialPeers := flag.String("peers", "ws://localhost:6001", "initial peers")
	flag.Parse()

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// handle
	http.HandleFunc("/deposit", handleDeposit)
	http.HandleFunc("/withdraw", handleWithdraw)

	// Start listening for incoming peer messages
	go handleMessages()

	// http serve
	log.Println("http server started on", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal("Could not serve http: ", err)
	}
}
