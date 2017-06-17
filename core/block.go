package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/exced/blockchain/crypto"
)

// Block is the basic component of a blockchain.
type Block struct {
	Mutex        *sync.Mutex
	Index        int64         `json:"index"`
	PreviousHash string        `json:"previousHash"`
	Timestamp    int64         `json:"timestamp"`
	Transactions *Transactions `json:"transactions"`
	Nonce        int           `json:"nonce"`
	Hash         string        `json:"hash"`
}

var genesisBlock = &Block{
	Mutex:        &sync.Mutex{},
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1496696844,
	Nonce:        0,
	Hash:         "bd125513b6f734f19d169f5a95e35765ccc5c438d937f728e5febd2322e3ddc4",
}

// ToHash hashes receiver block.
func (b *Block) ToHash() string {
	return crypto.ToHash(fmt.Sprintf("%d%s%d%v%d", b.Index, b.PreviousHash, b.Timestamp, b.Transactions, b.Nonce))
}

// IsValid retrieves the cryptographic validity between receiver block and given previous block.
func (b *Block) IsValid(pb *Block) bool {
	return b.Hash == b.ToHash() && b.Index == pb.Index+1 && b.PreviousHash == pb.Hash
}

// GenNext creates the next block of receiver block given hashed data.
func (b *Block) GenNext(transactions *Transactions) (nb *Block) {
	nb = &Block{
		Mutex:        &sync.Mutex{},
		Index:        b.Index + 1,
		PreviousHash: b.Hash,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		Nonce:        crypto.RandNonce(),
	}
	nb.Hash = nb.ToHash()
	return nb
}

// Mine looks for a nonce to satisfy given difficulty
func (b *Block) Mine(difficulty int) *Block {
	b.Nonce = crypto.RandNonce()
	return b
}
