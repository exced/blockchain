package crypto

import (
	"bytes"
	s "strings"
)

type PoW struct {
	Difficulty int
	HashRate   int
}

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
