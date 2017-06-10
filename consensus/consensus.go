package consensus

import (
	"log"

	"golang.org/x/net/websocket"

	"github.com/exced/blockchain/core"
)

type ConsensusAPI struct {
	consensus *consensus
}

func NewConsensusAPI() *ConsensusAPI {
	return &ConsensusAPI{consensus: &consensus{Peers: make(map[*Peer]bool), PoW: }}
}

func (api *ConsensusAPI) Connect(addr string) {
	ws, err := websocket.Dial(addr, "", addr)
	if err != nil {
		log.Println("dial to peer", err)
		continue
	}
	initConnection(ws)
}

// Consensus represents a set of peers which work together to build a valid blockchain.
type consensus struct {
	PoW   *pow.P         `json:"pow"`   // Proof of Work
	Peers map[*Peer]bool `json:"peers"` // connected peers
}

func initConnection(ws *websocket.Conn) {
	ws.Write(FetchMsg())
	go wsHandleP2P(ws)
}

func (c *Consensus) Register(conn *websocket.Conn, key string) {
	p := NewPeer(conn, key)
	c.Peers[p] = true

	msg := &PeerConnection{
		Peer:   p,
		Status: true,
	}

	// Send it out to every peer that is currently connected
	for peer := range c.Peers {
		err := peer.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			peer.Conn.Close()
			delete(c.Peers, peer)
		}
	}
}

// Broadcast given data to other peers
func (c *Consensus) Broadcast(self Peer, data interface{}) {

}

// Validate block
func (c *Consensus) Validate(b *core.Block) bool {
	return false
}
