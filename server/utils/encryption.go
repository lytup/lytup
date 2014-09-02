package utils

import (
	"crypto/sha256"

	"code.google.com/p/go.crypto/pbkdf2"
)

func HashPassword(pwd, salt string) []byte {
	return pbkdf2.Key([]byte(pwd), []byte(salt), 64000, sha256.Size, sha256.New)
}
