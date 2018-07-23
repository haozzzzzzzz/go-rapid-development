package service

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func (m *ServiceSource) generateBashTargetServer(shDir string) (err error) {
	bashTargetServerFilePath := fmt.Sprintf("%s/target_server.sh", shDir)
	err = ioutil.WriteFile(bashTargetServerFilePath, []byte(bashTargetServerFileText), os.ModePerm)
	if nil != err {
		logrus.Errorf("write bash target server file %q failed. %s.", bashTargetServerFilePath, err)
		return
	}
	return
}

var bashTargetServerFileText = `#!/usr/bin/env bash
# 此脚本在目标机器执行
stage=$1
if [ -z ${stage} ]
then
    stage="test"
fi

serviceName=$2
if [ -z ${stage} ]
then
    echo "serviceName is required"
    exit
fi

zipFile=deploy_${stage}_${serviceName}.zip

serviceDir=/data/apps/${serviceName}
logDir=/data/logs/${serviceName}
logPath=${logDir}/info.log

# 创建服务所在文件夹
if [ ! -e ${serviceDir} ]
then
    sudo mkdir -p ${serviceDir}
fi

if [ ! -e ${logDir} ]
then
    sudo mkdir -p ${logDir}
fi

if [ ! -e ${logPath} ]
then
    sudo touch ${logPath}
fi

# 清空日志
emptyLogCrontabSh=/etc/cron.daily/empty_${serviceName}_log
if [ ! -e ${emptyLogCrontabSh} ]
then
    sudo touch ${emptyLogCrontabSh}
    sudo chmod 777 ${emptyLogCrontabSh}
    echo "#!/bin/sh
    " | sudo tee ${emptyLogCrontabSh}
fi

echo "cat /dev/null > ${logPath}
" | sudo tee ${emptyLogCrontabSh}

# 定时清空/var/log/awslogs/.log的内容
emptyAwslogsLog=/etc/cron.daily/empty_awslogs_log
if [ ! -e ${emptyAwslogsLog} ]
then
    sudo touch ${emptyAwslogsLog}
    sudo chmod 777 ${emptyAwslogsLog}
    echo "#!bin/sh
systemctl stop awslogsd.service
cat /dev/null > /var/log/awslogs.log
systemctl start awslogsd.service
" | sudo tee ${emptyAwslogsLog}
fi

# 设置awslogs配置
awsLogsConf=/etc/awslogs/config/${serviceName}.conf
echo "[${logPath}]

encoding=utf_8

file = ${logPath}

buffer_duration = 5000

log_stream_name = {instance_id}

initial_position = end_of_file

log_group_name = ${stage}${logPath}" | sudo tee ${awsLogsConf}

sudo systemctl restart awslogsd.service
sudo systemctl enable awslogsd.service

# 创建bash文件
runSh=${serviceDir}/run.sh
if [ ! -e ${runSh} ]
then
    sudo touch ${runSh}
    sudo chmod 777 ${runSh}
fi

# 启动文件
sudo echo "#!/usr/bin/env bash
${serviceDir}/main > $logPath 2>&1
" > ${runSh}

# 创建service文件
userServiceName=user_service_${serviceName}
userServiceFilePath=/lib/systemd/system/${userServiceName}.service

# 更新服务配置
touch temp_user_service.service
echo "[Unit]
Description=${userServiceName}
After=network.target awslogsd.service

[Service]
Type=simple
WorkingDirectory=${serviceDir}
ExecStart=${runSh}
Restart=always

[Install]
WantedBy=multi-user.target
" > temp_user_service.service
sudo cp -u temp_user_service.service ${userServiceFilePath}
rm temp_user_service.service
sudo systemctl daemon-reload

# 停止服务
sudo systemctl stop ${userServiceName}

# 解压服务文件
sudo mv ${zipFile} ${serviceDir}
cd ${serviceDir}
sudo unzip -o ${zipFile}

# 启动服务
sudo systemctl start ${userServiceName}
sudo systemctl status ${userServiceName}
sudo systemctl enable ${userServiceName}

exit
`