package client

import (
	"github.com/exced/simple-blockchain/core"
)

type APIBlockchain struct {
	b *core.Blockchain
}

func NewAPIBlockchain(b *core.Blockchain) *APIBlockchain {
	return &APIBlockchain{b}
}

// Sync synchronizes the blockchain to have the latest version
func (api *APIBlockchain) Sync() {

}

type APIMiner struct {
	b *core.Blockchain
}

func NewPublicMinerAPI(b *core.Blockchain) *PublicMinerAPI {
	return &PublicMinerAPI{b}
}
