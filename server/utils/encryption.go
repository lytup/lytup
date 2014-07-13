package utils

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha256"
	"os"
)

var (
	SALT = []byte(os.Getenv("SALT"))
	KEY  = []byte(os.Getenv("KEY"))
)

func HashPassword(pwd []byte) []byte {
	return pbkdf2.Key(pwd, SALT, 4096, sha256.Size, sha256.New)
}
