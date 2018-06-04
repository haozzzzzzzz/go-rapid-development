package source

import (
	"fmt"
	"io/ioutil"

	"github.com/haozzzzzzzz/go-rapid-development/tools/api/com/proj"
	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateBashParams(shDir string) (err error) {
	bashParamsFilePath := fmt.Sprintf("%s/params.sh", shDir)
	err = ioutil.WriteFile(bashParamsFilePath, []byte(bashParamsFileText), proj.ProjectFileMode)
	if nil != err {
		logrus.Errorf("write params bash %q failed. %s.", bashParamsFilePath, err)
		return
	}
	return
}

var bashParamsFileText = `#!/usr/bin/env bash
stage=$1
if [ -z ${stage} ]
then
    stage="test"
fi

# 服务名
serviceName=
goRoot=/usr/local/go
goPath=

# 跳板机
jumpServerKey=
jumpServer=

# 目标机
targetServerKey=
targetServer=

if [[ -z ${serviceName} || -z ${goRoot} || -z ${goPath} || -z ${jumpServerKey} || -z ${jumpServer} || -z ${targetServerKey} || -z ${targetServer} ]]
then
    echo "lack of params"
    exit
fi
`
