package utils

import (
	"encoding/base64"
	"os/exec"
)

const (
	SIZE = "200x200"
)

func imageThumbnail(file string) (string, error) {
	cmd := exec.Command("gm", "convert", "-size", SIZE, file, "-resize",
		SIZE+"^", "+profile", "*", "-gravity", "Center", "-extent", SIZE, "-")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func videoThumbnail(file string) (string, error) {
	cmd := exec.Command("convert", file, "-resize", SIZE+"^", "-strip",
		"-gravity", "Center", "-extent", SIZE, "jpg:-")
	out, err := cmd.Output()

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func CreateThumbnail(file, typ string) (string, error) {
	if IsImage(typ) {
		return imageThumbnail(file)
	} else if IsVideo(typ) {
		file = file + "[10]" // 10th frame
		return videoThumbnail(file)
	} else {
		return "", nil
	}
}
