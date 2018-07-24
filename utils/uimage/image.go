package uimage

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

const TYPE_JPEG = "jpeg"
const TYPE_PNG = "png"
const TYPE_GIF = "gif"

func Decode(r io.Reader) (image.Image, string, error) {
	return image.Decode(r)
}

func Encode(imgType string, writer io.Writer, img image.Image) (err error) {
	switch imgType {
	case TYPE_JPEG:
		err = jpeg.Encode(writer, img, nil)
	case TYPE_PNG:
		err = png.Encode(writer, img)
	case TYPE_GIF:
		err = gif.Encode(writer, img, nil)
	default:
		err = errors.New("unknown image type")
	}
	return
}
