package core

import (
	"crypto/rsa"
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

func (t *Transaction) Cipher(privateKey *rsa.PrivateKey) string {

}

func (t *Transaction) Decipher() string {

}
