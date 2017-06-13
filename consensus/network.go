package consensus

// Network is a set of Peer that listen and serves each other
type Network struct {
	Peers []*Peer `json:"peers"`
}

func NewNetwork() *Network {
	return &Network{}
}

func (n *Network) RemoveByIndex(i int) {
	n.Peers = append(n.Peers[0:i], n.Peers[i+1:]...)
}

func (n *Network) Add(peer *Peer) {
	n.Peers = append(n.Peers, peer)
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
