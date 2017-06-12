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
	case "send":
		if flag.NArg() < 3 {
			log.Fatal("usage:\n\t send to amount\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Fatalf("given amount could not be parsed as int: %v: %v", amount, err)
		}
		send(*rsaFilePath, *httpPort, flag.Arg(1), amount)
	default:
		panic("command does not exist")
	}
}

// send cryptocurrency
func send(rsaFilePath string, httpPort int, to string, amount int64) {

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey

	// hash private key to get id
	hash := sha256.New()
	io.WriteString(hash, string(fmt.Sprintf("%v", rsaPrivateKey)))

	// transaction
	transaction := &core.Transaction{From: fmt.Sprintf("%x", hash.Sum(nil)), To: to, Amount: amount}

	// cipher transaction
	hash = sha256.New()
	io.WriteString(hash, string(fmt.Sprintf("%v", transaction)))
	sig, err := crypto.Sign(hash.Sum(nil), rsaPrivateKey)
	if err != nil {
		log.Fatalf("failed to sign hash %s: %v", hash.Sum(nil), err)
	}

	// prepare transaction message
	transactionMessage := &consensus.TransactionMessage{
		Signature:    sig,
		Hash:         hash.Sum(nil),
		RsaPublicKey: rsaPublicKey,
		Transaction:  transaction,
	}

	// send : HTTP POST
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(transactionMessage)
	resp, _ := http.Post(fmt.Sprintf("http://localhost:%d/send", httpPort), "application/json; charset=utf-8", b)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("Response:", string(body))
}
