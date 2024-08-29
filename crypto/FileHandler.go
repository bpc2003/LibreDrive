// crypto - handles file encryption and password hashing
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log"
)

// Encrypt - encrypts a buffer with a given key
func Encrypt(key, buf []byte) (ciphertext []byte) {
	key, _ = hex.DecodeString(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(buf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], buf)

	h := hmac.New(sha256.New, key)
	h.Write(buf)
	ciphertext = append(ciphertext, h.Sum(nil)...)
	return
}

// Decrypt - decrypts a buffer with a given key
func Decrypt(key, buf []byte) ([]byte, error) {
	key, _ = hex.DecodeString(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := buf[:aes.BlockSize]
	ciphertext := buf[aes.BlockSize : len(buf)-sha256.Size]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	h := hmac.New(sha256.New, key)
	h.Write(ciphertext)
	exp := buf[len(buf)-sha256.Size:]
	if hmac.Equal(h.Sum(nil), exp) {
		return ciphertext, nil
	} else {
		return nil, errors.New("Invalid key")
	}
}
