package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// 搜索
// 如果希望使用递归搜索的话，请使用filepath.Walk功能
func SearchFileNames(
	dir string,
	filter func(fileInfo os.FileInfo) bool,
	recurse bool,
) (result []string, err error) {
	result = make([]string, 0)
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
		}

		if filter == nil || filter(fileInfo) {
			result = append(result, fmt.Sprintf("%s/%s", dir, fileInfo.Name()))
		}
	}

	return
}

// 向上搜索
func SearchFilenameBackwardIterate(curDir string, name string, check func(curDir string, fileName string) (cont bool)) {
	if curDir == "" { // access top
		return
	}

	fileName := fmt.Sprintf("%s/%s", curDir, name)
	if PathExists(fileName) && check(curDir, fileName) == false {
		return
	}

	SearchFilenameBackwardIterate(ParentDir(curDir), name, check)

	return
}
