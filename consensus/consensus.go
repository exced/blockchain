package consensus

import (
	"time"

	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
	"github.com/gorilla/websocket"
)

// Consensus represents a "contract" between peers which work together to build a valid blockchain.
type Consensus struct {
	Next       time.Time     `json:"next"`       // time of next tick
	HashRate   time.Duration `json:"hashrate"`   // hashrate duration
	Difficulty int           `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
	Network    Network       `json:"network"`
}

func next(hashRate time.Duration) time.Time {
	return time.Now().Add(hashRate)
}

// NewConsensus returns new consensus
func NewConsensus() *Consensus {
	hashRate := time.Duration(600) * time.Second // 10 minutes
	return &Consensus{Next: next(hashRate), HashRate: hashRate, Difficulty: 4}
}

// Connect connects to peer address and await for its consensus response.
func Connect(url string) (*Consensus, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	c := &Consensus{}
	err = conn.ReadJSON(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Consensus) UpdateNext() {
	c.Next = next(c.HashRate)
}

// Broadcast given Message to all peers
func (c *Consensus) Broadcast(msg Message) {
	// Send it out to every peer that is currently connected
	for i, peer := range c.Network.Peers {
		err := peer.Conn.WriteJSON(msg)
		if err != nil {
			peer.Conn.Close()
			c.Network.RemoveByIndex(i)
		}
	}
}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return crypto.MatchHash(b.ToHash(), c.Difficulty)
}
