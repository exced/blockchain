package consensus

import (
	"crypto/rsa"

	"github.com/exced/blockchain/core"
)

// MessageType represents an enumeration for peers messages
type MessageType int

const (
	PeerStatus MessageType = iota // peer connect / disconnect
	Transaction
	Block
)

// Message represents msg communication between peers
type Message struct {
	Type    MessageType `json:"type"`
	Message interface{} `json:"message"`
}

// PeerStatusMessage is sent when a peer has connected
type PeerStatusMessage struct {
	PublicKey *rsa.PublicKey `json:"publickey"`
	Status    bool           `json:"status"` // true if connect, false if disconnect
}

// TransactionMessage is sent when a client do a send request
type TransactionMessage struct {
	Signature []byte         `json:"signature"`
	Hash      []byte         `json:"hash"`
	PublicKey *rsa.PublicKey `json:"publickey"`
}

// BlockMessage is sent to show up mined block to consensus
type BlockMessage struct {
	Block *core.Block `json:"block"` // mined block
	Peer  *Peer       `json:"peer"`  // miner validator
}
