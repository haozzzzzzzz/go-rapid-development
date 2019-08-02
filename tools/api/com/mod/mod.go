package mod

import (
	"bufio"
	"github.com/haozzzzzzzz/go-rapid-development/utils/file"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

func FindGoMod(curDir string) (
	found bool,
	modName string,
	modDir string,
) {
	var err error
	var goModFilename string
	curDir, errAbs := filepath.Abs(curDir)
	err = errAbs
	if nil != err {
		logrus.Errorf("get abs file dir failed. curDir: %s, error: %s.", curDir, err)
		return
	}

	file.SearchFilenameBackwardIterate(curDir, "go.mod", func(curDir string, fileName string) (cont bool) {
		cont = false
		goModFilename = fileName
		return
	})
	if goModFilename == "" {
		return
	}

	found = true

	f, err := os.Open(goModFilename)
	if nil != err {
		logrus.Errorf("open go.mod failed. filename: %s, error: %s.", goModFilename, err)
		return
	}
	defer func() {
		err = f.Close()
		if nil != err {
			logrus.Errorf("close file failed. filename: %s, error: %s.", goModFilename, err)
			return
		}
	}()

	reader := bufio.NewReader(f)
	moduleLine := ""
	for {
		bLine, isPrefix, err := reader.ReadLine()
		if nil != err {
			logrus.Errorf("read go mod file failed. %s", err)
			return
		}

		if isPrefix { // too long
			break
		}

		line := string(bLine)
		if line == "" {
			break
		}

		if strings.Contains(line, "module") {
			moduleLine = line
			break
		}
	}

	if moduleLine == "" {
		return
	}

	modName = strings.Replace(moduleLine, "module", "", 1)
	modName = strings.TrimSpace(modName)
	modDir = filepath.Dir(goModFilename)

	return
}
