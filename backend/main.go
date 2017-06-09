package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"gopkg.in/mgo.v2"

	"github.com/exced/blockchain/backend/api"
	"github.com/exced/blockchain/core"
)

var (
	blockchain   *core.Blockchain     // blockchain
	transactions []*core.Transaction  // pending transactions
	consensus    *consensus.Consensus // consensus
	upgrader     websocket.Upgrader
)

func main() {
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	dbAddr := flag.String("db", "mongodb://localhost/blockchain", "DB listen address")
	flag.Parse()

	// MongoDB Dial
	session, err := mgo.Dial(*dbAddr)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// init blockchain if no local one
	blockchain = core.NewBlockchain()

	// user storage and user API
	userAPI := api.NewUserAPI(session)

	// routes
	r := mux.NewRouter()

	// user logic
	r.HandleFunc("/login", userAPI.LoginUser).Methods("POST")
	r.HandleFunc("/signin", userAPI.SigninUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("POST")
	r.HandleFunc("/user/{id}", userAPI.PostUser).Methods("PUT")
	r.HandleFunc("/user/{id}", userAPI.DeleteUser).Methods("DELETE")

	// peer logic
	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// http serve
	log.Println("http server started on", *httpAddr)
	err = http.ListenAndServe(*httpAddr, r)
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

	// send Consensus to current logged in user

}
