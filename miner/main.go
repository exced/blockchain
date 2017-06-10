package main

import (
	"flag"
	"log"
	"net/http"
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
	httpAddr := flag.String("http", ":3000", "HTTP listen address")
	p2pAddr := flag.String("p2p", ":6000", "P2P listen address")

	consensusAPI := consensus.NewConsensusAPI()

	// Peer handshake protocol
	consensusAPI.Connect(*p2pAddr)

	// Blockchain fetch
	blockchain = consensusAPI.Fetch()

	// Handle peers
	http.HandleFunc("/ws", consensusAPI.handlePeerConnection)

	// Handle cli
	http.HandleFunc("/withdraw", WithConsensus(handleWithdraw, consensusAPI))
	http.HandleFunc("/deposit", WithConsensus(handleDeposit, consensusAPI))

	// http serve
	log.Println("http server started on", *httpAddr)
	err := http.ListenAndServe(*httpAddr, nil)
	if err != nil {
		log.Fatal("Could not serve http: ", err)
	}
}
