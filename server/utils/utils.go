package utils

import (
	"strings"
)

func IsImage(typ string) bool {
	return strings.HasPrefix(typ, "image")
}

func IsVideo(typ string) bool {
	return strings.HasPrefix(typ, "video")
}
