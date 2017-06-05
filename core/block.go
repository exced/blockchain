package core

import (
	"fmt"
	"time"

	"github.com/exced/simple-blockchain/crypto"
)

// Block is the basic component of a blockchain.
type Block struct {
	Index        int64  `json:"index"`
	PreviousHash string `json:"previousHash"`
	Timestamp    int64  `json:"timestamp"`
	Data         string `json:"data"`
	Hash         string `json:"hash"`
}

var genesisBlock = &Block{
	Index:        0,
	PreviousHash: "0",
	Timestamp:    1496696844,
	Data:         "Genesis",
	Hash:         "bd125513b6f734f19d169f5a95e35765ccc5c438d937f728e5febd2322e3ddc4",
}

// toHash hashes receiver block.
func (b *Block) toHash() string {
	return crypto.ToHash(fmt.Sprintf("%d%s%d%s", b.Index, b.PreviousHash, b.Timestamp, b.Data))
}

// isValid retrieves the cryptographic validity between receiver block and given previous block.
func (b *Block) isValid(pb *Block) bool {
	return b.Hash == b.toHash() && b.Index == pb.Index+1 && b.PreviousHash == pb.Hash
}

// genNext creates the next block of receiver block given data.
func (b *Block) genNext(data string) (nb *Block) {
	nb = &Block{
		Data:         data,
		PreviousHash: b.Hash,
		Index:        b.Index + 1,
		Timestamp:    time.Now().Unix(),
	}
	nb.Hash = nb.toHash()
	return nb
}
