package crypto

import "math/rand"

// RandNonce returns a random int between 0 and 10000
func RandNonce() int {
	return rand.Intn(10000)
}
