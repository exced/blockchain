package consensus

import "golang.org/x/net/websocket"

// Peer defines our peer user
type Peer struct {
	Conn *websocket.Conn `json:"conn"` // ws connection
	Key  string          `json:"key"`  // public key of the peer
}

func NewPeer(c *websocket.Conn, key string) *Peer {
	return &Peer{Conn: c, Key: key}
}

type PeerConnection struct {
	Peer   *Peer `json:"peer"`
	Status bool  `json:"status"` // true: connect, false:disconnect
}
