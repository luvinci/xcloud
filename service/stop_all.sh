#!/bin/bash

set -e

function stop_service() {
	pid=`ps -aux | grep -v grep | grep "service/bin/$1" | awk '{print $2}'`
	if [[ $pid != "" ]];then
		ps -aux | grep -v grep | grep "service/bin/$1" | awk '{print $2}' | xargs kill
		echo -e "\033[32m服务已关闭：$1\033[0m"
	else
		echo -e "\033[31m服务并未启动：$1\033[0m"
	fi
}

services="
apigw
account
dbproxy
upload
transfer
download
"

# 关闭service
for service in $services
do
	stop_service $service
done
