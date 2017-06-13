package consensus

import "time"

// Network is a set of Peer that listen and serves each other
type Network struct {
	Peers []*Peer `json:"peers"`
}

// NewNetwork retrieves a new empty network
func NewNetwork() *Network {
	return &Network{}
}

// RemoveByIndex remove a peer by its index position
func (n *Network) RemoveByIndex(i int) {
	n.Peers = append(n.Peers[0:i], n.Peers[i+1:]...)
}

// Add add a peer to the network
func (n *Network) Add(peer *Peer) {
	n.Peers = append(n.Peers, peer)
}

// Aggregate Peers answer
func (n *Network) Aggregate() []*Message {
	a := make([]*Message, len(n.Peers))
	for i, peer := range n.Peers {
		var msg *Message
		peer.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		// Read in a new message as JSON and map it to a Message object
		err := peer.Conn.ReadJSON(&msg)
		if err != nil {
			a[i] = msg
		}
	}
	return a
}

// Broadcast given Message to all peers
func (n *Network) Broadcast(msg Message) {
	// Send it out to every peer that is currently connected
	for i, peer := range n.Peers {
		err := peer.Conn.WriteJSON(msg)
		if err != nil {
			peer.Conn.Close()
			n.RemoveByIndex(i)
		}
	}
}
