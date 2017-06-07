package consensus

import (
	"github.com/exced/simple-blockchain/core"
	"github.com/exced/simple-blockchain/peer"
)

// Consensus represents a set of peers which work together to build a valid blockchain.
type Consensus map[*peer.Peer]bool

// Broadcast given data to other peers
func (c *Consensus) Broadcast(self peer.Peer, data interface{}) {

}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return false
}
