package consensus

import (
	"encoding/json"
	"time"

	"golang.org/x/net/websocket"

	"github.com/exced/blockchain/core"
)

// Consensus represents a set of peers which work together to build a valid blockchain.
type Consensus struct {
	Tick       time.Time                `json:"tick"`       // time of next tick
	HashRate   time.Duration            `json:"hashrate"`   // hashrate duration
	Difficulty int                      `json:"difficulty"` // number of 0 required at the beginning of the hash : Proof of Work
	Blockchain *Blockchain              `json:"blockchain"` // latest blockchain
	Peers      map[*websocket.Conn]bool `json:"peers"`      // connected peers
}

// NewConsensus returns new consensus
func NewConsensus() *Consensus {
	return &Consensus{}
}

// PeerConnectionMessage is a message sent when a peer has connected
type PeerConnectionMessage struct {
	Conn *websocket.Conn `json:"conn"`
}

// Connect connects to peer address and await for its consensus response
func Connect(addr string) (*Consensus, error) {
	ws, err := websocket.Dial(addr, "", addr)
	if err != nil {
		return nil, err
	}
	// get other peer consensus
	var c = &Consensus{}
	var msg []byte
	err := websocket.Message.Receive(ws, &msg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(msg, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// HandlePeerConnection sends Consensus to newly connected peer
func (c *Consensus) HandlePeerConnection(ws *websocket.Conn) error {
	consensus, err := json.Marshal(c)
	if err != nil {
		return err
	}
	ws.Write(consensus)
	msg := &PeerConnectionMessage{ws}
	// Broadcast New Peer
	for peer := range peers {
		peer.Conn.WriteJSON(msg)
	}
	return nil
}

func (c *Consensus) ListenAndServe() {
	var b *Block
	for {
		for x := range time.Tick(time.Until(c.Tick)) {
			b = Mine()
			if MatchHash()
		}
		// Broadcast Mined Block
		for peer := range peers {
			peer.Conn.WriteJSON(msg)
		}
		blocks, err := c.Synchronize()
		c.Validate(blocks)
		c.Tick += c.HashRate
	}

}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return false
}

// Mine mines block
func (c *Consensus) Mine() *Block {

}
