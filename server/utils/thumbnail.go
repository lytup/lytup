package utils

import (
	"github.com/gographics/imagick/imagick"
	"os"
	"path"
)

func CreateThumbnail(folPath, filePath, fileName, typ string) (bool, error) {
	// Only image or video
	if !IsImage(typ) && !IsVideo(typ) {
		return false, nil
	}

	if IsVideo(typ) {
		filePath += "[10]" // 10th frame
		fileName += ".jpg" // Save as image
	}

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(filePath)
	if err != nil {
		return false, err
	}

	tPath := path.Join(folPath, "t") // Thumbnail path
	err = os.MkdirAll(tPath, 0755)
	if err != nil {
		return false, err
	}

	mw.ResizeImage(200, 200, imagick.FILTER_UNDEFINED, 1)
	mw.WriteImage(path.Join(tPath, fileName))

	return true, nil
}
