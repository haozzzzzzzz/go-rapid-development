package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateBashBuild(shDir string) (err error) {
	bashBuildFilePath := fmt.Sprintf("%s/build.sh", shDir)
	err = ioutil.WriteFile(bashBuildFilePath, []byte(bashBuildFileText), os.ModePerm)
	if nil != err {
		logrus.Errorf("write bash build file %q failed. %s.", bashBuildFilePath, err)
		return
	}
	return
}

var bashBuildFileText = `#!/usr/bin/env bash
source params.sh

export GOROOT=${goRoot}
export GOPATH=${goPath}
export GOOS=linux
export GOARCH=amd64

api compile -p ../
go build -o ../stage/${stage}/main ../main.go
cd ../stage/${stage}/
zip -r deploy_${stage}_${serviceName}.zip main config
`
