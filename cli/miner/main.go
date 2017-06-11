package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
	c "github.com/exced/blockchain/consensus"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*pb.TransactionMessage
	consensus           *c.Consensus
)

func main() {
	p2pAddr := flag.String("p2p", "6000", "P2P listen address")
	rpcAddr := flag.String("http", "3000", "RPC address")
	rsaFilePath := flag.String("r", "../private.pem", "RSA key file")
	// blockchainFilePath := flag.String("b", "./blockchain", "Blockchain file")
	flag.Parse()

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
	if err != nil {
		log.Fatal("could not open rsa file", err)
	}
	rsaPublicKey := &rsaPrivateKey.PublicKey
	fmt.Println("rsapublickey ", rsaPublicKey)

	// Genesis Peer
	if *p2pAddr == "genesis" {
		consensus = c.NewConsensus()
	} else {
		consensus, err = c.Connect(*p2pAddr)
		if err != nil {
			log.Fatalf("could not connect to consensus %s : %v", *p2pAddr, err)
		}
	}

	// Handle peers connection
	// http.HandleFunc("/ws", websocket.Handler(consensus.HandlePeerConnection))
	http.HandleFunc("/blockchain", handleBlockchain)

	// Serve Network
	go consensus.ListenAndServe()

	// mine
	go func() {
		// b := c.Blockchain.Mine(c.Difficulty)
		// if crypto.MatchHash(b.ToHash(), c.Difficulty) {

		// }
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
}

func handleBlockchain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(consensus.Blockchain)
}

type server struct{}

func (server) Send(ctx context.Context, t *pb.TransactionMessage) (*pb.Response, error) {
	fmt.Println("SEND ", *t)
	pendingTransactions = append(pendingTransactions, t)
	return &pb.Response{Success: true, Msg: "successfully added to pending transactions"}, nil
}
