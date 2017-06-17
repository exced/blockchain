package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"crypto/rsa"

	"github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/core"
	"github.com/exced/blockchain/crypto"
)

func main() {
	httpPort := flag.Int("p", 3000, "HTTP port send address")
	rsaFilePath := flag.String("i", "./private.pem", "RSA key file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t send to amount \n\t gen")
	}

	switch flag.Arg(0) {
	case "gen":
		crypto.GenRsaFile(*rsaGenFilePath)
	case "withdraw":
		if flag.NArg() < 3 {
			log.Fatal("usage:\n\t withdraw to amount\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Fatalf("given amount could not be parsed as int: %v: %v", amount, err)
		}
		withdraw(*httpPort, *rsaFilePath, flag.Arg(1), amount)
	case "deposit":
		if flag.NArg() < 3 {
			log.Fatal("usage:\n\t deposit from amount\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Fatalf("given amount could not be parsed as int: %v: %v", amount, err)
		}
		deposit(*httpPort, *rsaFilePath, flag.Arg(1), amount)
	default:
		panic("command does not exist")
	}
}

// deposit cryptocurrency
func deposit(httpPort int, rsaFilePath string, from string, amount int64) {
	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey

	// hash private key to get id
	hash := sha256.New()
	io.WriteString(hash, string(fmt.Sprintf("%v", rsaPrivateKey)))

	addr := fmt.Sprintf("http://localhost:%d/transaction", httpPort)
	send(addr, from, fmt.Sprintf("%x", hash.Sum(nil)), amount, rsaPrivateKey, rsaPublicKey)
}

// withdraw cryptocurrency
func withdraw(httpPort int, rsaFilePath string, to string, amount int64) {
	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey

	// hash private key to get id
	hash := sha256.New()
	io.WriteString(hash, string(fmt.Sprintf("%v", rsaPrivateKey)))

	addr := fmt.Sprintf("http://localhost:%d/transaction", httpPort)
	send(addr, fmt.Sprintf("%x", hash.Sum(nil)), to, amount, rsaPrivateKey, rsaPublicKey)
}

func send(addr, from, to string, amount int64, rsaPrivateKey *rsa.PrivateKey, rsaPublicKey *rsa.PublicKey) {
	// transaction
	transaction := &core.Transaction{From: from, To: to, Amount: amount}
	transactionMessage := consensus.NewTransactionMessage(transaction, rsaPrivateKey, rsaPublicKey)

	// HTTP POST
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(transactionMessage)
	resp, err := http.Post(addr, "application/json; charset=utf-8", b)
	if err != nil {
		log.Fatalf("failed to post to %s : %v", addr, err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response:", string(body))
}
