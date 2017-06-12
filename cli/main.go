package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"io"
	"log"
	"strconv"

	"github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

func main() {
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	rsaFilePath := flag.String("r", "./private.pem", "RSA key file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t send from to amount \n\t gen")
	}

	switch flag.Arg(0) {
	case "gen":
		crypto.GenRsaFile(*rsaGenFilePath)
	case "send":
		if flag.NArg() < 4 {
			log.Fatal("usage:\n\t send from to amount\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(3), 10, 64)
		if err != nil {
			log.Fatalf("given amount could not be parsed as int: %v: %v", amount, err)
		}
		send(*rsaFilePath, *httpAddr, flag.Arg(1), flag.Arg(2), amount)
	default:
		panic("command does not exist")
	}
}

// send cryptocurrency
func send(rsaFilePath, rpcAddr, from, to string, amount int64) {

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey

	// transaction
	transaction := &core.Transaction{From: *rsaPublicKey, To: to, Amount: amount}
	transactionString, err := json.Marshal(transaction)
	if err != nil {
		log.Fatalf("could not marshal transaction: %#v: %v", transaction, err)
	}

	hash := sha256.New()
	io.WriteString(hash, string(transactionString))
	sig, err := crypto.Sign(hash.Sum(nil), rsaPrivateKey)
	if err != nil {
		log.Fatalf("failed to sign hash %s: %v", hash.Sum(nil), err)
	}

	rsaPublicKeyBytes, err := crypto.GetBytes(rsaPublicKey)
	if err != nil {
		log.Fatalf("failed to get bytes of %v : %v", rsaPublicKey, err)
	}
	transactionMessage := &consensus.TransactionMessage{Signature: sig, Hash: hash.Sum(nil), Rsapublickey: rsaPublicKeyBytes}

	// send
}
