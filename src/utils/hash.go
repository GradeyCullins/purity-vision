package utils

import (
	"crypto/sha256"
	"fmt"
)

// Hash returns a sha256 hash of the string argument in hex format.
func Hash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
