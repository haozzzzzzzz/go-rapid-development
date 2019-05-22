#!/usr/bin/env bash
stage=$1
if [ -z ${stage} ]
then
    stage="test"
fi

# 服务名
serviceName=test_compile
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
