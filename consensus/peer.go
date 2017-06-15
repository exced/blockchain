package consensus

import (
	"github.com/gorilla/websocket"
)

// Peer is a wrapper of a websocket connection.
type Peer struct {
	Address string          `json:"address"`
	Conn    *websocket.Conn `json:"conn"`
}
