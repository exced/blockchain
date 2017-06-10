package core

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"time"

	"github.com/exced/blockchain/crypto"
)

// Block is the basic component of a blockchain.
type Block struct {
	Index        int64          `json:"index"`
	PreviousHash string         `json:"previousHash"`
	Timestamp    int64          `json:"timestamp"`
	Data         string         `json:"data"`
	Transactions []*Transaction `json:"transactions"`
	Nonce        int            `json:"nonce"`
	Hash         string         `json:"hash"`
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
	return crypto.ToHash(fmt.Sprintf("%d%s%d%s%v%d", b.Index, b.PreviousHash, b.Timestamp, b.Data, b.Transactions, b.Nonce))
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

// NextBlock looks for a nonce to satisfy given PoW
func (b *Block) NextBlock(p *crypto.PoW) *Block {
	for b.Nonce = rand.Intn(10000); !p.MatchHash(b.ToHash()); {

	}
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

func verify(sig, hash []byte, publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, sig)
}

	err = verify(sig, hash.Sum(nil), rsaPublicKey)
	if err != nil {
		log.Fatal(err.Error())
	}
