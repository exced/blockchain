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
	Transactions: NewTransactions(),
	Nonce:        0,
	Hash:         "c02c463d7c5559d90a6b90facf87de3f451aed75d26e8dadc778d4e140f59beb",
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
func (b *Block) GenNext(t *Transactions) (nb *Block) {
	nb = &Block{
		Mutex:        &sync.Mutex{},
		Index:        b.Index + 1,
		PreviousHash: b.Hash,
		Timestamp:    time.Now().Unix(),
		Transactions: t,
		Nonce:        crypto.RandNonce(),
	}
	nb.Hash = nb.ToHash()
	return nb
}

// Link received block to given previous block : set the previousHash field to previous Hash
func (b *Block) Link(pb *Block) {
	b.PreviousHash = pb.Hash
}

// Mine generates a rand nonce for the block
func (b *Block) Mine() {
	b.Nonce = crypto.RandNonce()
}
