package consensus

import (
	"crypto/rsa"

	"golang.org/x/net/websocket"
)

// Peer represents our miner
type Peer struct {
	Conn      *websocket.Conn `json:"conn"`
	PublicKey *rsa.PublicKey  `json:"publickey"`
}

// NewPeer initializes and return a new Peer
func NewPeer(conn *websocket.Conn, publickey *rsa.PublicKey) *Peer {
	return &Peer{Conn: conn, PublicKey: publickey}
}
