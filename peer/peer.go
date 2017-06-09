package peer

import "golang.org/x/net/websocket"

// PeerConnection defines our peer connection message object
type PeerConnection struct {
	Conn   *websocket.Conn `json:"conn"`   // ws connection
	Status bool            `json:"status"` // connect: true, disconnect: false
}
