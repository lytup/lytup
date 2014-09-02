package utils

import (
	"strings"

	"github.com/dchest/uniuri"
)

func IsImage(typ string) bool {
	return strings.HasPrefix(typ, "image")
}

func IsVideo(typ string) bool {
	return strings.HasPrefix(typ, "video")
}

// func Salt(n uint8) ([]byte, error) {
// 	salt := make([]byte, n)
// 	if _, err := rand.Read(salt); err != nil {
// 		return nil, err
// 	}
// 	return salt, nil
// }

func RandomString(n int) string {
	return uniuri.NewLen(n)
}
