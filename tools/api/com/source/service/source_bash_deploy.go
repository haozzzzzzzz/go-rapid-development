package service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateBashDeploy(shDir string) (err error) {
	bashDeployFilePath := fmt.Sprintf("%s/deploy.sh", shDir)
	err = ioutil.WriteFile(bashDeployFilePath, []byte(bashDeployFileText), os.ModePerm)
	if nil != err {
		logrus.Errorf("write bash deploy file %q failed. %s.", bashDeployFilePath, err)
		return
	}

	return
}

var bashDeployFileText = `#!/usr/bin/env bash
source params.sh

read -p "Input essh password:" -s password

echo ""
# 复制文件到跳板机
cd ../stage/${stage}
echo "[begin] copy package to remote"
sshpass -p ${password} scp -o StrictHostKeyChecking=no -v ./jump_${stage}_${serviceName}.tar.gz ${jumpServer}:~/luohao/
echo "[finish] copy package to remote"

cd ../../sh/
echo "[begin] deploy"
# 在跳板机执行操作
sshpass -p ${password} ssh -o StrictHostKeyChecking=no ${jumpServer} 'bash -s' < "jump_server.sh" ${stage} ${serviceName} ${targetServerKey} ${targetServer}
echo "[finish] deploy"`
