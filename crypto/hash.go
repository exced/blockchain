package crypto

import (
	"crypto/sha256"
	"fmt"
)

// ToHash hashes given string using SHA256.
func ToHash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
