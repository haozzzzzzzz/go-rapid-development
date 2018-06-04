package source

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ApiProjectSource) generateBashDeploy(shDir string) (err error) {
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

# 复制文件到跳板机
cd ../stage/${stage}
scp -i ${jumpServerKey} ./deploy_${stage}_${serviceName}.zip ${jumpServer}:~/luohao/

cd ../../sh/
scp -i ${jumpServerKey} target_server.sh ${jumpServer}:~/luohao/

# 在跳板机执行操作
ssh -i ${jumpServerKey} ${jumpServer} 'bash -s' < "jump_server.sh" ${stage} ${serviceName} ${targetServerKey} ${targetServer}
`
