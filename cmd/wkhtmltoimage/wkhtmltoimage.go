package wkhtmltoimage

import (
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/sirupsen/logrus"
)

// 输入html文件地址，生成图片文件

func WkHtmlToImage(
	inputPath string,
	outputPath string,
	format string,
) (err error) {
	_, err = cmd.RunCommand("~/", "wkhtmltoimage", "--format", format, inputPath, outputPath)
	if err != nil {
		logrus.Errorf("run cmd wkhtmtoimage to generate picture failed. error: %s", err)
		return
	}
	return
}
