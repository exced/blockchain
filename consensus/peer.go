package consensus

import (
	"log"

	"github.com/gorilla/websocket"
)

// Peer Listen and Serve a websocket connection
type Peer struct {
	Conn *websocket.Conn `json:"conn"`
}

func (p *Peer) ListenAndServe() {
	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		// handling message
		switch msg.Type {
		case Transaction:
			log.Println("msg Transaction")
		case Block:
			log.Println("msg Block")
		}
	}
}
