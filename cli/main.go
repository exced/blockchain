package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
)

func main() {
	peerAddr := flag.String("http", ":3000", "Peer address")
	rsaFilePath := flag.String("i", "./private.pem", "RSA key file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t \"send\"\n \n gen")
	}

	switch flag.Arg(0) {
	case "gen":
		genRsaFile(*rsaGenFilePath)
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
		rsaPrivateKey, err := openRsaFile(*rsaFilePath)
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
		sig, err := sign(hash.Sum(nil), rsaPrivateKey)
		if err != nil {
			log.Fatal(err)
		}
		hashBytes, err := getBytes(hash)
		if err != nil {
			log.Fatal(err)
		}
		rsaPublicKeyBytes, err := getBytes(&rsaPrivateKey.PublicKey)
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

func sign(hash []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
}

// Generate RSA Private Key and store it in a file
func genRsaFile(path string) error {
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaPrivateKey),
		})
	return ioutil.WriteFile(path, pemdata, 0644)
}

func openRsaFile(path string) (*rsa.PrivateKey, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(f)
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
