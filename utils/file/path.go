package file

import (
	"os"
)

func PathExists(path string) (exists bool) {
	_, err := os.Stat(path)
	if nil != err {
		return
	}

	exists = true
	return
}
