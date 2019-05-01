package file

import (
	"os"
	"strings"

	"github.com/haozzzzzzzz/go-rapid-development/utils/str"
)

func PathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	if nil != err {
		return
	}

	exists = true
	return
}

func ParentDir(dir string) string {
	return str.SubString(dir, 0, strings.LastIndex(dir, "/"))
}
