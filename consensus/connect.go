package consensus

import (
	"github.com/exced/blockchain/core"
	"github.com/gorilla/websocket"
)

// DialResponse is sent to a peer who has recrently connected. It gives min datas to start working.
type DialResponse struct {
	Conn       *websocket.Conn  `json:"conn"`
	Consensus  *Consensus       `json:"consensus"`
	Network    *Network         `json:"network"`
	Blockchain *core.Blockchain `json:"blockchain"`
}

// Connect connects to peer address and await for its dial response
func Connect(url string) (*DialResponse, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	r := &DialResponse{}
	err = conn.ReadJSON(r)
	if err != nil {
		return nil, err
	}
	r.Conn = conn
	return r, nil
}
