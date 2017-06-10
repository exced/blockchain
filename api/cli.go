package api

import (
	"crypto/rsa"

	"github.com/exced/blockchain/core"
)

type CliAPI struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewCliAPI(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *CliAPI {
	return &CliAPI{privateKey, publicKey}
}

func Deposit(from string, currency core.Currency, amount int64) {

}

func Withdraw(to string, currency core.Currency, amount int64) {

}
