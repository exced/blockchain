package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/exced/blockchain/consensus"
)

type ConsensusAPI struct {
	Consensus *consensus.Consensus
}

func NewConsensusAPI() *ConsensusAPI {
	return &ConsensusAPI{Consensus: consensus.NewConsensusAPI()}
}

// handlePeerConnection broadcast newly connected or disconnected peer to peers
func (api *ConsensusAPI) handlePeerConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var v struct {
		Key string `json:"key"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err = decoder.Decode(&v)
	if err != nil {

		return
	}

	// Register our new peer
	consensus.Register(ws, v.Key)

	// send Consensus to current logged in user
	ws.WriteJSON(consensus)
}
