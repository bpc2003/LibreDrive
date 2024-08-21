package crypto

import (
  "crypto/sha256"
  "math/rand/v2"
)

func GeneratePassword(password string, rounds int) (string, string) {
  var salt string
  for i := 0; i < rounds; i++ {
    salt += string(rand.Int() % 128)
  }

  h := sha256.Sum256([]byte(password + salt))
  password = string(h[:])
  return password, salt
}

func ComparePassword(password, salt, hashed string) bool {
  h := sha256.Sum256([]byte(password + salt))
  return string(h[:]) == hashed
}
