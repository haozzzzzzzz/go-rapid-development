package file

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
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

// if file witch specified extension exists in directory
func DirHasExtFile(dir string, ext string) (hasExt bool, err error) {
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		logrus.Errorf("read directory %s failed. error: %s.", dir, err)
		return
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		fileExt := filepath.Ext(fileInfo.Name())
		if fileExt == ext {
			hasExt = true
			break
		}
	}

	return
}
