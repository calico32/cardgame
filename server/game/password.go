package game

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	"golang.org/x/crypto/sha3"
)

const saltLength = 16

func hashPassword(password string) string {
	var salt [saltLength]byte
	n, err := rand.Read(salt[:])
	if n != saltLength || err != nil {
		panic("failed to generate salt: " + err.Error())
	}

	h := sha3.New512()
	h.Write(salt[:])
	h.Write([]byte(password))
	return fmt.Sprintf("%x%x", salt, h.Sum(nil))
}

func checkPassword(saltedHash, password string) bool {
	saltString, hash := saltedHash[:saltLength*2], saltedHash[saltLength*2:]
	salt, err := hex.DecodeString(saltString)
	if err != nil {
		panic("failed to decode salt: " + err.Error())
	}

	h := sha3.New512()
	h.Write(salt)
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil)) == hash
}
