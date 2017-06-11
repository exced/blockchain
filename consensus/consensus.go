package consensus

import (
	"encoding/json"
	"time"

	"golang.org/x/net/websocket"

	"crypto/rsa"

	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

// Consensus represents a set of peers which work together to build a valid blockchain.
// Each peer is a slave for the network and a master for client
type Consensus struct {
	Tick       time.Time        `json:"tick"`       // time of next tick
	HashRate   time.Duration    `json:"hashrate"`   // hashrate duration
	Difficulty int              `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
	Blockchain *core.Blockchain `json:"blockchain"` // latest blockchain
	Peers      map[*Peer]bool   `json:"peers"`      // connected peers
}

// NewConsensus returns new consensus
func NewConsensus() *Consensus {
	return &Consensus{}
}

// Connect connects to peer address and await for its consensus response.
func Connect(addr string) (*Consensus, error) {
	ws, err := websocket.Dial(addr, "", addr)
	if err != nil {
		return nil, err
	}
	// get other peer consensus
	var c = &Consensus{}
	var msg []byte
	err = websocket.Message.Receive(ws, &msg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msg, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Connect notify myself to network
func (c *Consensus) Connect(publickey *rsa.PublicKey) {
	msg, _ := json.Marshal(&PeerStatusMessage{PublicKey: publickey, Status: true})
	// Broadcast New Peer
	for peer := range c.Peers {
		peer.Conn.Write(msg)
	}
}

// HandlePeerConnection sends Consensus to newly connected peer
func (c *Consensus) HandlePeerConnection(ws *websocket.Conn) error {
	consensus, err := json.Marshal(c)
	if err != nil {
		return err
	}
	ws.Write(consensus)
	return nil
}

// ListenAndServe serves network by mining and handling consensus queries
func (c *Consensus) ListenAndServe() {
	broadcast := make(chan Message)
	go c.handleBroadcast(broadcast)

	for {
		// time
		go func() {
			for range time.Tick(time.Until(c.Tick)) {
			}
			// c.Tick += c.HashRate
		}()
	}

}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return crypto.MatchHash(b.ToHash(), c.Difficulty)
}

func (c *Consensus) handleBroadcast(broadcast <-chan Message) {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		mmsg, _ := json.Marshal(msg)
		// Send it out to every peer that is currently connected
		for peer := range c.Peers {
			_, err := peer.Conn.Write(mmsg)
			if err != nil {
				peer.Conn.Close()
				delete(c.Peers, peer)
			}
		}
	}
}
