package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"io"
	"log"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
	"github.com/exced/blockchain/crypto"
)

func main() {
	peerAddr := flag.String("http", ":3000", "Peer address")
	rsaFilePath := flag.String("r", "./private.pem", "RSA key file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t \"send\"\n \n gen")
	}

	switch flag.Arg(0) {
	case "gen":
		crypto.GenRsaFile(*rsaGenFilePath)
	case "send":
		if flag.NArg() < 3 {
			log.Fatal("usage:\n\t \"send key amount\"\n")
		}
		// args
		amount, err := strconv.ParseInt(flag.Arg(2), 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		// rsa key
		rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
		if err != nil {
			log.Fatal(err)
		}
		// transaction
		if err != nil {
			log.Fatal(err.Error())
		}
		// gRPC
		conn, err := grpc.Dial(*peerAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect to %s: %v", *peerAddr, err)
		}
		defer conn.Close()
		client := pb.NewPeerClient(conn)

		// transaction
		transaction := &pb.Transaction{To: flag.Arg(0), Amount: amount}
		transactionString, err := json.Marshal(transaction)
		if err != nil {
			log.Fatal(err)
		}
		hash := sha256.New()
		io.WriteString(hash, string(transactionString))
		sig, err := crypto.Sign(hash.Sum(nil), rsaPrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		hashBytes, err := crypto.GetBytes(hash)
		if err != nil {
			log.Fatal(err)
		}
		rsaPublicKeyBytes, err := crypto.GetBytes(&rsaPrivateKey.PublicKey)
		if err != nil {
			log.Fatal(err)
		}
		transactionMessage := &pb.TransactionMessage{Signature: sig, Hash: hashBytes, Rsapublickey: rsaPublicKeyBytes}

		// send
		res, err := client.Send(context.Background(), transactionMessage)
		if err != nil {
			log.Fatalf("could not send %d to %s: %v", amount, flag.Arg(0), err)
		}
		log.Println(res)
	default:
		panic("command does not exist")
	}
}
