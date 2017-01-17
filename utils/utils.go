package utils

import (
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/scrypt"
)

const (
	scryptN    = 1048576
	scryptR    = 8
	scryptP    = 1
	passLength = 32
)

func NewUUID() string {
	return uuid.NewV4().String()
}

func HashWithNewSalt(source string) (string, string, error) {
	salt := NewUUID()
	dk, err := scrypt.Key([]byte(source), []byte(salt), scryptN, scryptR, scryptP, passLength)
	return string(dk), salt, err
}
