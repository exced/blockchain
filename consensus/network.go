package consensus

// Network is a set of Peer that listen and serves each other
type Network struct {
	Peers []*Peer `json:"peers"`
}

func (n *Network) RemoveByIndex(i int) {
	n.Peers = append(n.Peers[0:i], n.Peers[i+1:]...)
}

func (n *Network) Add(peer *Peer) {
	n.Peers = append(n.Peers, peer)
}

func (n *Network) ListenAndServe() {
	for _, peer := range n.Peers {
		go peer.ListenAndServe()
	}
}
