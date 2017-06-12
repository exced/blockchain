package consensus

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"crypto/rsa"

	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

// Consensus represents a set of peers which work together to build a valid blockchain.
// Each peer is a slave for the network and a master for client
type Consensus struct {
	Tick       time.Time                `json:"tick"`       // time of next tick
	HashRate   time.Duration            `json:"hashrate"`   // hashrate duration
	Difficulty int                      `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
	Peers      map[*websocket.Conn]bool `json:"peers"`      // connected peers
}

// NewConsensus returns new consensus
func NewConsensus() *Consensus {
	hashRate := time.Duration(600) * time.Second // 10 minutes
	tick := time.Now().Add(hashRate)
	return &Consensus{Tick: tick, HashRate: hashRate, Difficulty: 4}
}

// Connect connects to peer address and await for its consensus response.
func Connect(url string) (*Consensus, error) {
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var c = &Consensus{}
	err = json.Unmarshal(msg, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Connect notify myself to network
func (c *Consensus) Connect(publickey *rsa.PublicKey) {
	// Broadcast New Peer
	for peer := range c.Peers {
		peer.WriteJSON(&PeerStatusMessage{PublicKey: publickey, Status: true})
	}
}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return crypto.MatchHash(b.ToHash(), c.Difficulty)
}

// Broadcast given Message to all peers
func (c *Consensus) Broadcast(msg Message) {
	// Send it out to every peer that is currently connected
	for peer := range c.Peers {
		err := peer.WriteJSON(msg)
		if err != nil {
			peer.Close()
			delete(c.Peers, peer)
		}
	}
}
