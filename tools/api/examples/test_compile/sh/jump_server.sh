#!/usr/bin/env bash
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

cd luohao

jumpZipFile=jump_${stage}_${serviceName}.tar.gz
tar -xzvf ${jumpZipFile}

zipFile=deploy_${stage}_${serviceName}.zip
targetServerSh=target_server.sh

if [ ! -e ${zipFile} ]
then
    echo "no target zip file ${zipFile} exists"
    exit
fi

# 拷贝文件到远程主机
scp -o StrictHostKeyChecking=no -i ${targetServerKey} ${zipFile} ${targetServer}:~/
ssh -o StrictHostKeyChecking=no -i ${targetServerKey} ${targetServer} 'bash -s' < ${targetServerSh} ${stage} ${serviceName}

# 删除本地临时文件
rm ${jumpZipFile}
rm ${zipFile}
rm ${targetServerSh}
