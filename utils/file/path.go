package file

import (
	"os"

	"github.com/sirupsen/logrus"
)

func PathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	if nil != err {
		logrus.Errorf("path: %p not exists. \n%s.", path, err)
		return
	}

	exists = true
	return
}
