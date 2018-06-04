package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateBashJumpServer(shDir string) (err error) {
	bashJumpServerFilePath := fmt.Sprintf("%s/jump_server.sh", shDir)
	err = ioutil.WriteFile(bashJumpServerFilePath, []byte(bashDeployFileText), os.ModePerm)
	if nil != err {
		logrus.Errorf("write bash jump server file %q failed. %s.", bashJumpServerFilePath, err)
		return
	}
	return
}

var bashJumpServerFileText = `#!/usr/bin/env bash
# 此脚本在跳板机执行
stage=$1
serviceName=$2
targetServerKey=$3
targetServer=$4

if [[ -z ${stage} || -z ${serviceName} || -z ${targetServerKey} || -z ${targetServer} ]]
then
    echo "lack of params"
    exit
fi

zipFile=deploy_${stage}_${serviceName}.zip
targetServerSh=target_server.sh

cd luohao
if [ ! -e ${zipFile} ]
then
    echo "no target zip file ${zipFile} exists"
    exit
fi

# 拷贝文件到远程主机
scp -i ${targetServerKey} ${zipFile} ${targetServer}:~/
ssh -i ${targetServerKey} ${targetServer} 'bash -s' < ${targetServerSh} ${stage} ${serviceName}

# 删除本地临时文件
rm ${zipFile}
rm ${targetServerSh}
`
