package main

import (
	xlog "github.com/luvinci/x-logrus"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/sirupsen/logrus"
	"time"
	cfg "xcloud/config"
	"xcloud/service/download/proto"
	"xcloud/service/download/route"
	"xcloud/service/download/rpc"
)

var (
	consulAddr = cfg.Viper.GetString("consul.addr")
	downloadPort =  cfg.Viper.GetString("app.download.port")
)

func startDownloadRpcService() {
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.download"),
		micro.RegisterTTL(10*time.Second),     // 超时时间，避免注册中心没有主动删除已失去心跳的节点
		micro.RegisterInterval(5*time.Second), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		// micro.Flags(common.CustomFlags...),
	)

	service.Init()

	_ = proto.RegisterDownloadServiceHandler(service.Server(), new(rpc.Download))
	if err := service.Run(); err != nil {
		logrus.Error(err)
	}
}

func startDownloadApiService() {
	r := route.Router()
	err := r.Run(":" + downloadPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	xlog.Init()
	// 上传 api 服务
	go startDownloadApiService()
	// 上传 rpc 服务
	startDownloadRpcService()
}
