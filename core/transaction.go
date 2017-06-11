package core

import (
	"crypto/rsa"

	"github.com/exced/blockchain/crypto"
)

// Transaction represents a money transaction between 2 users
type Transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int64  `json:"amount"`
}

func NewTransaction(from, to string, a int64) *Transaction {
	return &Transaction{from, to, a}
}

func (t *Transaction) Verify(sig, hash []byte, rsaPublicKey *rsa.PublicKey) error {
	return crypto.Verify(sig, hash, rsaPublicKey)
}
