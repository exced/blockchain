package consensus

import (
	"github.com/exced/blockchain/core"
	"github.com/gorilla/websocket"
)

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
	conn.WriteJSON(peer)
	r := &DialResponse{}
	err = conn.ReadJSON(r)
	if err != nil {
		return nil, nil, err
	}
	return conn, r, nil
}
