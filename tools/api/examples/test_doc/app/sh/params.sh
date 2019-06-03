#!/usr/bin/env bash
stage=$1
if [ -z ${stage} ]
then
    stage="test"
fi

# 服务名
serviceName=test_doc_app
goRoot=/usr/local/go
goPath=/Users/hao/Documents/Projects/XunLei/video_buddy_service

# 跳板机
jumpServerKey=xx
jumpServer=xx

# 目标机
targetServerKey=xx
targetServer=xx

if [[ -z ${serviceName} || -z ${goRoot} || -z ${goPath} || -z ${jumpServerKey} || -z ${jumpServer} || -z ${targetServerKey} || -z ${targetServer} ]]
then
    echo "lack of params"
    exit
fi
