package wkhtmltoimage

import (
	"fmt"
	"testing"
)

func TestWkHtmlToImage(t *testing.T) {
	//err := WkHtmlToImage("/home/hao/Projects/tool/wkhtmltoimage/sample.html", "./sample.png", "png")
	err := WkHtmlToImage(&Params{
		InputUri:   "http://api.videobuddy.vid007.com",
		OutputPath: "./sample.png",
		Format:     "png",
		Width:      0,
		Height:     0,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("finish")
}
