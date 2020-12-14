package wkhtmltoimage

import (
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
)

type Params struct {
	WorkingDir string `json:"working_dir"`
	ExecName   string `json:"exec_name"`
	InputUri   string `json:"input_uri"`
	OutputPath string `json:"output_path"`
	Format     string `json:"format"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

func (m *Params) Logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"working_dir": m.WorkingDir,
		"exec_name":   m.ExecName,
		"input_uri":   m.InputUri,
		"output_path": m.OutputPath,
		"format":      m.Format,
		"width":       m.Width,
		"height":      m.Height,
	})
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
	if params.WorkingDir == "" {
		params.WorkingDir = "./"
	}

	if params.ExecName == "" {
		params.ExecName = "wkhtmltoimage"
	}

	logger := params.Logger()

	_, err = cmd.RunCommand(
		"./",
		"wkhtmltoimage",
		args...,
	)
	if err != nil {
		logger.Errorf("run cmd wkhtmtoimage to generate picture failed. error: %s", err)
		return
	}

	return
}
