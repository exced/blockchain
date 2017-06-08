package model

// Currency represents the currency we use in our blockchain.
type Currency int

const (
	EXC Currency = iota
	BTC
	ETH
)
