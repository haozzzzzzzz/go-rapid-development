package uimage

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"io"

	"github.com/sirupsen/logrus"
	_ "golang.org/x/image/webp" // webp只有decode，没有encode

	"golang.org/x/image/bmp"
)

const TYPE_JPEG = "jpeg"
const TYPE_PNG = "png"
const TYPE_GIF = "gif"
const TYPE_BMP = "bmp"

func Decode(r io.Reader) (img image.Image, imgType string, err error) {
	img, imgType, err = image.Decode(r)
	if nil != err {
		logrus.Errorf("decode image failed. %s.", err)
		return
	}
	return
}

// encode后可能会改变图片类型
func Encode(imgType string, writer io.Writer, img image.Image) (newImgType string, err error) {
	newImgType = imgType
	switch imgType {
	case TYPE_JPEG:
		err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 50})
	case TYPE_PNG:
		err = png.Encode(writer, img)
	case TYPE_GIF:
		err = gif.Encode(writer, img, nil)
	case TYPE_BMP:
		err = bmp.Encode(writer, img)
	default: // webp
		newImgType = TYPE_JPEG
		logrus.Warnf("encode unsupported encode image type: %s to jpeg", imgType)
		err = jpeg.Encode(writer, img, &jpeg.Options{Quality: 50}) // 转成png
	}
	return
}
