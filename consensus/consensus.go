package consensus

import (
	"encoding/json"
	"time"

	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

// Consensus represents a "contract" between peers which work together to build a valid blockchain.
type Consensus struct {
	Next       time.Time     `json:"next"`       // time of next tick
	HashRate   time.Duration `json:"hashrate"`   // hashrate duration
	Difficulty int           `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
}

// UpdateNext update the time of the next Tick adding hashRate to Now()
func (c *Consensus) UpdateNext() {
	c.Next = time.Now().Add(c.HashRate)
}

// NewConsensus returns new consensus
func NewConsensus() *Consensus {
	hashRate := time.Duration(3) * time.Second // 10 minutes
	return &Consensus{Next: time.Now().Add(hashRate), HashRate: hashRate, Difficulty: 4}
}

// Validate given block according to given responses and received consensus
func (c *Consensus) Validate(block *core.Block, responses []*Message) bool {
	noise := 0
	valid := 0
	for _, resp := range responses {
		if resp.Type == Block {
			var blockMsg *BlockMessage
			err := json.Unmarshal(resp.Message, &blockMsg)
			if (err != nil) || (blockMsg.Block != block) {
				noise++
				continue
			}
			if blockMsg.Signature {
				valid++
			}
		}
	}
	return (2*valid > len(responses)-noise) && crypto.MatchHash(block.ToHash(), c.Difficulty)
}
