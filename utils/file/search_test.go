package file

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestSearchFileNames(t *testing.T) {
	SearchFileNames("/Users/hao/Documents/Projects/XunLei/video_buddy_service/src/service/common/model/../", func(fileInfo os.FileInfo) bool {
		if strings.HasSuffix(fileInfo.Name(), "r.go") {
			fmt.Println(fileInfo.Name())
		}
		return true
	}, true)
}
