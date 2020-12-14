package wkhtmltoimage

import (
	"fmt"
	"testing"
)

func TestWkHtmlToImage(t *testing.T) {
	err := WkHtmlToImage("/home/hao/Projects/tool/wkhtmltoimage/sample.html", "sample.png", "png")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("finish")
}
