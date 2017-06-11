package crypto

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	s "strings"
)

// ToHash hashes given string using SHA256.
func ToHash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

// MatchHash test if given hash has given difficulty number of "0" at the beginning
func MatchHash(hash string, difficulty int) bool {
	var buffer bytes.Buffer
	for i := 0; i < difficulty; i++ {
		buffer.WriteString("0")
	}
	return s.HasPrefix(hash, buffer.String())
}
