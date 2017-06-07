package crypto

import (
	"crypto/sha256"
	"fmt"
)

// ToHash hashes given string using SHA256.
func ToHash(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

// Hasher represents a hashable item i.e. an object that has a ToHash function
type Hasher interface {
	ToHash() string
}
