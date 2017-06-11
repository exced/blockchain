package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"golang.org/x/net/websocket"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*Transaction
	consensus           *c.Consensus
)

func main() {
	p2pAddr := flag.String("p2p", ":6000", "P2P listen address")
	rpcAddr := flag.String("http", ":3000", "RPC address")
	rsaFilePath := flag.String("r", "./private.pem", "RSA key file")
	blockchainFilePath := flag.String("b", "./blockchain", "Blockchain file")
	rsaGenFilePath := flag.String("o", "./private.pem", "RSA key generated file")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("usage:\n\t \"send\"\n \n gen \n mine")
	}

	switch flag.Arg(0) {
	case "gen":
		crypto.GenRsaFile(*rsaGenFilePath)
	case "mine":
		// rsa key
		rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
		if err != nil {
			log.Fatal(err)
		}

		// Genesis Peer
		if *p2pAddr == "genesis" {
			consensus = c.NewConsensus()
		} else {
			consensus, err = c.Connect(*p2pAddr)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Handle peers connection
		http.HandleFunc("/ws", websocket.Handler(consensus.HandlePeerConnection))
		http.HandleFunc("/blockchain", handleBlockchain)

		// Serve Network
		go consensus.ListenAndServe()

		// mine
		go func() {
			b := c.Blockchain.Mine(c.Difficulty)
			if crypto.MatchHash(b.ToHash(), c.Difficulty) {

			}
			// Broadcast Mined Block
			// broadcast <- &Message{Type: BlockMessage, Message: &BlockMessage{Block: b}}
		}()

		// gRPC Server
		log.Printf("listening to port %s", *rpcAddr)
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *rpcAddr))
		if err != nil {
			log.Fatalf("could not listen to port %s: %v", *rpcAddr, err)
		}
		s := grpc.NewServer()
		pb.RegisterPeerServer(s, server{})
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("could not serve: %v", err)
		}
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
		// gRPC Cli
		conn, err := grpc.Dial(*rpcAddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect to %s: %v", *rpcAddr, err)
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

func handleBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(consensus.Blockchain)
}

type server struct{}

func (server) Send(ctx context.Context, t *pb.TransactionMessage) (*pb.Response, error) {

	pendingTransactions = append(pendingTransactions, t)
	consensus.Broadcast()
}
