// 这个工具有个劣势：对于html、css的特性支持不足
package wkhtmltoimage

import (
	"context"
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
)

type LoadErrorHandleType string

const (
	LoadErrorHandleAbort  LoadErrorHandleType = "abort"
	LoadErrorHandleIgnore LoadErrorHandleType = "ignore"
	LoadErrorHandleSkip   LoadErrorHandleType = "skip"
)

type LoadMediaErrorHandleType string

const (
	LoadMediaErrorHandleAbort  LoadMediaErrorHandleType = "abort"
	LoadMediaErrorHandleIgnore LoadMediaErrorHandleType = "ignore"
	LoadMediaErrorHandlerSkip  LoadMediaErrorHandleType = "skip"
)

type LogLevelType string

const (
	LogLevelTypeNone  LogLevelType = "none"
	LogLevelTypeError LogLevelType = "error"
	LogLevelTypeWarn  LogLevelType = "warn"
	LogLevelTypeInfo  LogLevelType = "info"
)

// man wkhtmltoimage
type Params struct {
	WorkingDir             string                   `json:"working_dir"`
	ExecName               string                   `json:"exec_name"`
	InputUri               string                   `json:"input_uri"`
	OutputPath             string                   `json:"output_path"`
	Format                 string                   `json:"format"`
	Width                  int                      `json:"width"`  // 像素宽度。
	Height                 int                      `json:"height"` // 像素高度。
	NoDebugJavascript      bool                     `json:"no_debug_javascript"`
	DisableJavascript      bool                     `json:"disable_javascript"`
	Quality                int                      `json:"quality"`                   // 图片质量。0～100
	EnableLocalFileAccess  bool                     `json:"enable_local_file_access"`  // 允许访问本地文件
	LoadErrorHandling      LoadErrorHandleType      `json:"load_error_handling"`       // 页面加载失败处理。
	LoadMediaErrorHandling LoadMediaErrorHandleType `json:"load_media_error_handling"` // 媒体加载失败处理
	LogLevel               LogLevelType             `json:"log_level"`
	Allow                  string                   `json:"allow"` // 允许程序加载的本地目录下的文件
}

func (m *Params) Logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"working_dir":               m.WorkingDir,
		"exec_name":                 m.ExecName,
		"input_uri":                 m.InputUri,
		"output_path":               m.OutputPath,
		"format":                    m.Format,
		"width":                     m.Width,
		"height":                    m.Height,
		"no_debug_javascript":       m.NoDebugJavascript,
		"disable_javascript":        m.DisableJavascript,
		"quality":                   m.Quality,
		"enable_local_file_access":  m.EnableLocalFileAccess,
		"load_error_handling":       m.LoadErrorHandling,
		"load_media_error_handling": m.LoadMediaErrorHandling,
		"log_level":                 m.LogLevel,
		"allow":                     m.Allow,
	})
}

// 输入html uri，生成图片文件
func WkHtmlToImage(
	ctx context.Context,
	params *Params,
) (err error) {
	args := make([]string, 0)
	args = append(args, "--format", params.Format)

	if params.Width > 0 {
		args = append(args, "--width", to.String(params.Width))
	}

	if params.Height > 0 {
		args = append(args, "--height", to.String(params.Height))
	}

	if params.NoDebugJavascript {
		args = append(args, "--no-debug-javascript")
	}

	if params.DisableJavascript {
		args = append(args, "--disable-javascript")
	}

	if params.Quality > 0 {
		args = append(args, "--quality", to.String(params.Quality))
	}

	if params.EnableLocalFileAccess {
		args = append(args, "--enable-local-file-access")
	}

	if params.LoadErrorHandling != "" {
		args = append(args, "--load-error-handling", string(params.LoadErrorHandling))
	}

	if params.LoadMediaErrorHandling != "" {
		args = append(args, "--load-media-error-handling", string(params.LoadMediaErrorHandling))
	}

	if params.LogLevel != "" {
		args = append(args, "--log-level", to.String(params.LogLevel))
	}

	if params.Allow != "" {
		args = append(args, "--allow", params.Allow)
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

	_, err = cmd.RunCommandCtx(
		ctx,
		params.WorkingDir,
		params.ExecName,
		args...,
	)
	if err != nil {
		logger.Errorf("run cmd wkhtmtoimage to generate picture failed. error: %s", err)
		return
	}

	return
}
