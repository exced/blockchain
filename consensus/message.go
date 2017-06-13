package consensus

import (
	"crypto/rsa"
	"encoding/json"

	"github.com/exced/blockchain/core"
)

// MessageType represents an enumeration for peers messages
type MessageType int

// PeerStatus, Transaction, Block represents Type for corresponding message
const (
	PeerStatus = iota
	Transaction
	Block
)

// Message represents msg communication between peers
type Message struct {
	Type    MessageType     `json:"type"`
	Message json.RawMessage `json:"message"`
}

// PeerStatusMessage is sent when a peer connect or disconnect
type PeerStatusMessage struct {
	Peer   *Peer `json:"peer"`
	Status bool  `json:"status"`
}

// TransactionMessage is sent when a client do a send request
type TransactionMessage struct {
	Signature    []byte            `json:"signature"`
	Hash         []byte            `json:"hash"`
	RsaPublicKey *rsa.PublicKey    `json:"rsapublickey"`
	Transaction  *core.Transaction `json:"transaction"`
}

// BlockMessage is sent to show up mined block to consensus
type BlockMessage struct {
	Block  *core.Block `json:"block"`  // mined block
	From   string      `json:"from"`   // from ID to know who rewards
	Status bool        `json:"status"` // true: accepted, false: rejected
}
