package service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateBashBuild(shDir string) (err error) {
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
curDir=` + "`pwd`" + `
cd ../stage/${stage}/
#cp ${goPath}/bin/logfmt ./ # copy log format tool
zip -r deploy_${stage}_${serviceName}.zip main config
echo "package deploy_${stage}_${serviceName}.zip finish"

cp ${curDir}/target_server.sh ./
tar -czvf jump_${stage}_${serviceName}.tar.gz deploy_${stage}_${serviceName}.zip target_server.sh
echo "package jump_${stage}_${serviceName}.tar.gz finish"
`
