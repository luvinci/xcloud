package main

import (
	xlog "github.com/luvinci/x-logrus"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/sirupsen/logrus"
	"os"
	"time"
	cfg "xcloud/config"
	"xcloud/mq"
	"xcloud/service/upload/proto"
	"xcloud/service/upload/route"
	"xcloud/service/upload/rpc"
)

var (
	consulAddr = cfg.Viper.GetString("consul.addr")
	uploadPort =  cfg.Viper.GetString("app.upload.port")
)

func startUploadRpcService() {
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.upload"),
		micro.RegisterTTL(10*time.Second),     // 超时时间，避免注册中心没有主动删除已失去心跳的节点
		micro.RegisterInterval(5*time.Second), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		// micro.Flags(common.CustomFlags...),
	)

	// TODO: 检查是否指定dbhost（原）
	service.Init()

	_ = proto.RegisterUploadServiceHandler(service.Server(), new(rpc.Upload))
	if err := service.Run(); err != nil {
		logrus.Error(err)
	}
}

func startUploadApiService() {
	r := route.Router()
	err := r.Run(":" + uploadPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	xlog.Init()
	mq.Init()
	temp := cfg.Viper.GetString("file.TempLocalRootDir")
	chunk := cfg.Viper.GetString("file.TempChunkRootDir")
	_ = os.MkdirAll(temp, os.ModePerm)
	_ = os.MkdirAll(chunk, os.ModePerm)

	// 上传 api 服务
	go startUploadApiService()
	// 上传 rpc 服务
	 startUploadRpcService()
}
