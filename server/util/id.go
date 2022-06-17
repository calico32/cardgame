package util

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid"
	"golang.org/x/crypto/sha3"
)

func RoomCode() string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"
	return "r" + gonanoid.MustGenerate(alphabet, 6)
}

func IdFrom(prefix string, text string) string {
	h := fmt.Sprintf("%x", sha3.Sum256([]byte(text)))
	return fmt.Sprintf("%s_%s", prefix, h[:8])
}

func LongIdFrom(prefix string, text string) string {
	h := fmt.Sprintf("%x", sha3.Sum256([]byte(text)))
	return fmt.Sprintf("%s_%s", prefix, h)
}
