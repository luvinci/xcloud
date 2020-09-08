#!/bin/bash

set -e

# ROOT_DIR
ROOT_DIR=/web/xcloud

cd $ROOT_DIR

# 打包静态资源（以防配置被修改）
mkdir -p ${ROOT_DIR}/assets && go-bindata-assetfs -pkg assets -o ${ROOT_DIR}/assets/asset.go static/... config/config.ini

# 编译service
function build_service() {
    echo -e "\033[35m开始编译：service/$1/main.go\033[0m"
    go build -o service/bin/$1 service/$1/main.go
    resbin=`ls service/bin/ | grep $1`
    echo -e "\033[32m编译完成：service/bin/${resbin}\033[0m"
}

services="
apigw
account
dbproxy
upload
transfer
download
"

# 执行编译每个service
mkdir -p service/bin/ && rm -f service/bin/*
for service in $services
do
    build_service $service
done
