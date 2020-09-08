#!/bin/bash

set -e

# ROOT_DIR
ROOT_DIR=/web/xcloud

cd $ROOT_DIR

# 创建运行日志目录
logpath=/web/logs/xcloud
mkdir -p $logpath

services="
apigw
account
dbproxy
upload
transfer
download
"

# 启动service
function start_service() {
    nohup ./service/bin/$1 >> $logpath/$1.log 2>&1 &
    if [[ $? == 0 ]];then
    	echo -e "\033[32m服务启动成功：service/bin/$1\033[0m"
   	else
   		echo -e "\033[31m服务启动失败：service/bin/$1\033[0m"
   	fi
}

# 执行启动service
for service in $services
do
    start_service $service
done
