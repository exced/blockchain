package consensus

import (
	"encoding/json"

	"github.com/exced/blockchain/core"
	"github.com/gorilla/websocket"
)

// Peer is a wrapper of a websocket connection.
type Peer struct {
	Address string          `json:"address"`
	Conn    *websocket.Conn `json:"conn"`
}

// NewPeer returns a new peer instance with given connect address
func NewPeer(address string) *Peer {
	return &Peer{Address: address}
}

// DialResponse is sent to a peer who has recrently connected. It gives min datas to start working.
type DialResponse struct {
	Consensus  *Consensus       `json:"consensus"`
	Network    *Network         `json:"network"`
	Blockchain *core.Blockchain `json:"blockchain"`
}

// Connect connects to peer address and await for its dial response
func Connect(url string, peer *Peer) (*websocket.Conn, *DialResponse, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, nil, err
	}
	msg, _ := json.Marshal(&PeerConnMessage{Peer: peer})
	conn.WriteJSON(Message{Type: PeerConn, Message: msg})
	r := &DialResponse{}
	err = conn.ReadJSON(r)
	if err != nil {
		return nil, nil, err
	}
	return conn, r, nil
}
