package auth

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

const (
	KeyTime    = 1
	KeyMemory  = 64 * 1024
	KeyThreads = 4
	KeyLen     = 32
	SaltLen    = 32
)

func Salt() ([]byte, error) {
	b := make([]byte, SaltLen)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Key(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, KeyTime, KeyMemory, KeyThreads, KeyLen)
}

func Encode(password []byte) (key, salt []byte, err error) {
	salt, err = Salt()
	if err != nil {
		return
	}
	key = Key(password, salt)
	return
}
