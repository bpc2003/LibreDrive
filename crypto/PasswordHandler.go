package crypto

import (
  "crypto/sha256"
  "fmt"
  "math/rand/v2"
)

func GeneratePassword(password string, rounds int) (string, string) {
  var salt string
  for i := 0; i < rounds; i++ {
    ch := string(rand.Int() % 256)
    salt += ch
    h := sha256.Sum256([]byte(password + ch))
    password = string(h[:])
  }

  return password, salt
}

func ComparePassword(password, salt, hashed string) bool {
  var h [sha256.Size]byte
  for _, r := range salt {
    h = sha256.Sum256([]byte(password + fmt.Sprintf("%c", r)))
    password = string(h[:])
  }

  return string(h[:]) == hashed
}
