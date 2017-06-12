package consensus

import (
	"github.com/gorilla/websocket"
)

// Network represent connected Peers abstraction
type Network struct {
	Peers map[*websocket.Conn]string
}

func NewNetwork() *Network {
	return &Network{make(map[*websocket.Conn]string)}
}

// Connect notify myself to network
func (n *Network) Connect(conn *websocket.Conn) {
	for peer := range n.Peers {
		peer.WriteJSON(&PeerStatusMessage{Conn: conn, Status: true})
	}
}

// Broadcast given Message to all peers
func (n *Network) Broadcast(msg Message) {
	// Send it out to every peer that is currently connected
	for peer := range n.Peers {
		err := peer.WriteJSON(msg)
		if err != nil {
			peer.Close()
			delete(n.Peers, peer)
		}
	}
}
