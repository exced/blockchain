package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/exced/blockchain/cli/api"
	"github.com/exced/blockchain/crypto"
)

var (
	pendingTransactions []*Transaction
)

// WithConsensus API
func WithConsensus(h http.Handler, c *consensus.consensusAPI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r)
	})
}

func main() {
	rpcAddr := flag.String("http", ":3000", "RPC address")
	rsaFilePath := flag.String("r", "./private.pem", "RSA key file")
	blockchainFilePath := flag.String("b", "./blockchain", "Blockchain file")
	p2pAddr := flag.String("p2p", ":6000", "P2P listen address")

	// rsa key
	rsaPrivateKey, err := crypto.OpenRsaFile(*rsaFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// API
	consensusAPI := consensus.NewConsensusAPI()

	// Peer handshake protocol
	consensusAPI.Connect(*p2pAddr)

	// Blockchain fetch
	blockchain = consensusAPI.Fetch()

	// Handle peers
	http.HandleFunc("/ws", consensusAPI.handlePeerConnection)

	// gRPC Server
	log.Println("listening to port %s", *rpcAddr)
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

	// Mine
	go Mine()
}

type server struct{}

func (server) Send(ctx context.Context, t *pb.TransactionMessage) (*pb.Response, error) {
	pendingTransactions = append(pendingTransactions, t)
	consensusAPI.Broadcast()
}

func Mine() {
	for {
		consensusAPI.Present()
	}
}
