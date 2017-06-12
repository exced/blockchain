package core

import (
	"fmt"
	"time"

	"github.com/exced/blockchain/crypto"
)

// Block is the basic component of a blockchain.
type Block struct {
	Index        int64          `json:"index"`
	PreviousHash string         `json:"previousHash"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	Nonce        int            `json:"nonce"`
	Hash         string         `json:"hash"`
}

var genesisBlock = &Block{
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

// isValid retrieves the cryptographic validity between receiver block and given previous block.
func (b *Block) isValid(pb *Block) bool {
	return b.Hash == b.ToHash() && b.Index == pb.Index+1 && b.PreviousHash == pb.Hash
}

// genNext creates the next block of receiver block given hashed data.
func (b *Block) genNext(transactions []*Transaction) (nb *Block) {
	nb = &Block{
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

// areTransactionsValid tests if given transactions are valid with current transactions chain
func (b *Block) areTransactionsValid(transactions []*Transaction) bool {
	balances := make(map[string]int64) // key - amount
	for _, t := range b.Transactions {
		balances[t.From] -= t.Amount
		balances[t.To] += t.Amount
	}
	for _, t := range transactions {
		balances[t.From] -= t.Amount
		if balances[t.From] < 0 {
			return false
		}
		balances[t.To] += t.Amount
	}
	return true
}
