package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
	"github.com/exced/blockchain/crypto"
)

func main() {
	rpcAddr := flag.String("http", "localhost:3000", "RPC address")
	rsaFilePath := flag.String("r", "./private.pem", "RSA key file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t send key amount \n\t gen")
	}

	switch flag.Arg(0) {
	case "gen":
		crypto.GenRsaFile(*rsaGenFilePath)
	case "send":
		if flag.NArg() < 3 {
			log.Fatal("usage:\n\t send key amount\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Fatalf("given amount could not be parsed as int: %v: %v", amount, err)
		}
		send(*rsaFilePath, *rpcAddr, flag.Arg(1), amount)
	default:
		panic("command does not exist")
	}
}

// send cryptocurrency
func send(rsaFilePath, rpcAddr, to string, amount int64) {
	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	// gRPC Cli
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", rpcAddr, err)
	}
	defer conn.Close()
	client := pb.NewPeerClient(conn)

	// transaction
	transaction := &pb.Transaction{To: to, Amount: amount}
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

	rsaPublicKeyBytes, err := crypto.GetBytes(&rsaPrivateKey.PublicKey)
	if err != nil {
		log.Fatalf("failed to get bytes of %v : %v", &rsaPrivateKey.PublicKey, err)
	}
	transactionMessage := &pb.TransactionMessage{Signature: sig, Hash: hash.Sum(nil), Rsapublickey: rsaPublicKeyBytes}

	// send
	res, err := client.Send(context.Background(), transactionMessage)
	if err != nil {
		log.Fatalf("could not send %d to %s: %v", amount, flag.Arg(0), err)
	}
	fmt.Println("RES ", res)
}
