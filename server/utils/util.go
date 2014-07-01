package utils

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"crypto/sha256"
)

var (
	SALT = []byte("Aiv5chie")
	KEY  = []byte("uSee4gee")
)

func HashPassword(pwd []byte) []byte {
	return pbkdf2.Key(pwd, SALT, 4096, sha256.Size, sha256.New)
}
