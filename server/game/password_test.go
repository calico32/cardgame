package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashGenerate(t *testing.T) {
	password := "correct horse battery staple"
	hash := hashPassword(password)
	assert.NotEmpty(t, hash, "hash is empty")
}

func TestHashCheck(t *testing.T) {
	password := "correct horse battery staple"
	hash := hashPassword(password)
	assert.True(t, checkPassword(hash, password), "password check failed")
}

func TestRoomPassword(t *testing.T) {
	r := Room{}
	password := "correct horse battery staple"
	r.SetPassword(password)
	assert.True(t, r.TryPassword(password), "room password check failed")
}
