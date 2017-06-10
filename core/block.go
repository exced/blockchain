package core

import (
	"fmt"
	"math/rand"
	"time"

	"crypto/rsa"

	"github.com/exced/blockchain/crypto"
)

// Block is the basic component of a blockchain.
type Block struct {
	Index        int64                           `json:"index"`
	PreviousHash string                          `json:"previousHash"`
	Timestamp    int64                           `json:"timestamp"`
	Data         string                          `json:"data"`
	Transactions map[*rsa.PublicKey]*Transaction `json:"transactions"`
	Nonce        int                             `json:"nonce"`
	Hash         string                          `json:"hash"`
}

var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1496696844,
	Data:         "Genesis",
	Nonce:        0,
	Hash:         "bd125513b6f734f19d169f5a95e35765ccc5c438d937f728e5febd2322e3ddc4",
}

// ToHash hashes receiver block.
func (b *Block) ToHash() string {
	return crypto.ToHash(fmt.Sprintf("%d%s%d%s", b.Index, b.PreviousHash, b.Timestamp, b.Data))
}

// isValid retrieves the cryptographic validity between receiver block and given previous block.
func (b *Block) isValid(pb *Block) bool {
	return b.Hash == b.ToHash() && b.Index == pb.Index+1 && b.PreviousHash == pb.Hash
}

// genNext creates the next block of receiver block given hashed data.
func (b *Block) genNext(data string) (nb *Block) {
	nb = &Block{
		Data:         data,
		PreviousHash: b.Hash,
		Index:        b.Index + 1,
		Timestamp:    time.Now().Unix(),
	}
	nb.Hash = nb.ToHash()
	return nb
}

func (b *Block) NextBlock(p *crypto.PoW) *Block {
	for b.Nonce = rand.Intn(10000); !p.MatchHash(b.ToHash()); {

	}
	return b
}

func (b *Block) areTransactionsValid() bool {
	for _, t := range b.Transactions {

	}
}
