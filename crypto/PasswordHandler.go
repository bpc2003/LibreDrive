package crypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"math/rand/v2"
)

// GeneratePassword - generates a password with a salt
func GeneratePassword(password string, rounds int) (string, string) {
	var salt string
	for i := 0; i < rounds; i++ {
		ch := string(rand.Int() % 256)
		salt += ch
		h := sha256.Sum256([]byte(password + ch))
		password = string(h[:])
	}
	h := sha512.Sum512([]byte(password))

	return fmt.Sprintf("%x", h), salt
}

// ComparePassword - compares a password and salt with a hashed password
func ComparePassword(password, salt, hashed string) bool {
	var h [sha256.Size]byte
	for _, r := range salt {
		h = sha256.Sum256([]byte(password + fmt.Sprintf("%c", r)))
		password = string(h[:])
	}

	return fmt.Sprintf("%x", sha512.Sum512(h[:])) == hashed
}
