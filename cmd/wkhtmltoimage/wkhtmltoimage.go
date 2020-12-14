package wkhtmltoimage

import (
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
)

type Params struct {
	InputUri   string `json:"input_uri"`
	OutputPath string `json:"output_path"`
	Format     string `json:"format"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

// 输入html uri，生成图片文件
func WkHtmlToImage(params *Params) (err error) {
	args := make([]string, 0)
	args = append(args, "--format", params.Format)

	if params.Width > 0 {
		args = append(args, "--width", to.String(params.Width))
	}

	if params.Height > 0 {
		args = append(args, "--height", to.String(params.Height))
	}

	// final
	args = append(args, params.InputUri, params.OutputPath)
	_, err = cmd.RunCommand(
		"./",
		"wkhtmltoimage",
		args...,
	)
	if err != nil {
		logrus.Errorf("run cmd wkhtmtoimage to generate picture failed. error: %s", err)
		return
	}

	return
}
