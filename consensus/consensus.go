package consensus

import (
	"time"

	"github.com/gorilla/websocket"

	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

// Consensus represents a set of peers which work together to build a valid blockchain.
// Each peer is a slave for the network and a master for client
type Consensus struct {
	Tick       time.Time     `json:"tick"`       // time of next tick
	HashRate   time.Duration `json:"hashrate"`   // hashrate duration
	Difficulty int           `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
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
	c := &Consensus{}
	err = conn.ReadJSON(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return crypto.MatchHash(b.ToHash(), c.Difficulty)
}
