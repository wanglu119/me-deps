package util

import (
	"fmt"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func CompareHashAndPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GenerateRandPassword(passwordLen int) string {
	b := make([]byte, passwordLen)
	var randVal uint32
	for i := 0; i < passwordLen; i++ {
		byteIdx := i % 4
		if byteIdx == 0 {
			randVal = rand.Uint32()
		}
		b[i] = byte((randVal >> (8 * uint(byteIdx)) & 0xFF))
	}
	tmpStr := fmt.Sprintf("%x", b)
	if len(tmpStr) > passwordLen {
		return tmpStr[:passwordLen]
	}
	return tmpStr
}
