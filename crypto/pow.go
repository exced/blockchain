package crypto

import (
	"bytes"
	s "strings"
)

type PoW struct {
	Difficulty int // number of first 0 recquired
	HashRate   int // time in seconds
}

var (
	initPoW = &PoW{Difficulty: 4, HashRate: 5000}
)

func NewPoW(difficulty, hashrate int) *PoW {
	return &PoW{Difficulty: difficulty, HashRate: hashrate}
}

func (p *PoW) MatchHash(hash string) bool {
	var buffer bytes.Buffer
	for i := 0; i < p.Difficulty; i++ {
		buffer.WriteString("0")
	}
	return s.HasPrefix(hash, buffer.String())
}
