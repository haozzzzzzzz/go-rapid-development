package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// 搜索
func SearchFileNames(
	dir string,
	filter func(fileInfo os.FileInfo) bool,
	recurse bool,
) (result []string, err error) {
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		logrus.Errorf("read directory %q files failed.", dir)
		return
	}

	for _, fileInfo := range files {

		if fileInfo.IsDir() && recurse {
			subFiles, errSearch := SearchFileNames(fmt.Sprintf("%s/%s", dir, fileInfo.Name()), filter, recurse)
			if nil != errSearch {
				err = errSearch
				logrus.Errorf("recurse search file failed. \n%s.", err)
				return
			}

			result = append(result, subFiles...)
			continue
		}

		if filter == nil || filter(fileInfo) {
			result = append(result, fmt.Sprintf("%s/%s", dir, fileInfo.Name()))
		}
	}

	return
}
