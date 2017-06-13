package consensus

import (
	"crypto/rsa"

	"github.com/exced/blockchain/core"
)

// MessageType represents an enumeration for peers messages
type MessageType int

// PeerStatus, Transaction, Block represents Type for corresponding message
const (
	Transaction = iota
	Block
)

// Message represents msg communication between peers
type Message struct {
	Type    MessageType `json:"type"`
	Message interface{} `json:"message"`
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
	Block *core.Block `json:"block"` // mined block
	From  string      `json:"from"`  // from ID to know who rewards
}
