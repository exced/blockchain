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
	PeerConn = iota
	Transaction
	Block
	BlockPoW
)

// Message represents msg communication between peers
type Message struct {
	Type    MessageType     `json:"type"`
	Message json.RawMessage `json:"message"`
}

// TransactionMessage is sent when a client do a send request
type TransactionMessage struct {
	Signature    []byte            `json:"signature"`
	Hash         []byte            `json:"hash"`
	RsaPublicKey *rsa.PublicKey    `json:"rsapublickey"`
	Transaction  *core.Transaction `json:"transaction"`
}

// PeerConnMessage is sent when a peer connect to the network
type PeerConnMessage struct {
	Peer *Peer `json:"peer"`
}

// BlockMessage is sent to show up mined block to consensus
type BlockMessage struct {
	Block     *core.Block `json:"block"`     // mined block
	From      string      `json:"from"`      // from ID to know who rewards
	Signature bool        `json:"signature"` // signature of peer who validated this block
}

// BlockPoWMessage is sent to notify other peers that a block has been accepted by the consensus
type BlockPoWMessage struct {
	Block      *core.Block `json:"block"`      // mined block
	From       string      `json:"from"`       // from ID to know who rewards
	Signatures []*Message  `json:"signatures"` // signature of peer who validated this block
}
