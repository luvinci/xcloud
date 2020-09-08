package main

import (
	xlog "github.com/luvinci/x-logrus"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/sirupsen/logrus"
	"time"
	cfg "xcloud/config"
	"xcloud/service/dbproxy/proto"
	"xcloud/service/dbproxy/rpc"
)

var (
	consulAddr = cfg.Viper.GetString("consul.addr")
)

// 启动 db rpc 服务
func startDbRpcService() {
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.dbproxy"),
		micro.RegisterTTL(10*time.Second), // 超时时间，避免注册中心没有主动删除已失去心跳的节点
		micro.RegisterInterval(5*time.Second), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		// micro.Flags(common.CustomFlags...),
	)

	// TODO: 检查是否指定dbhost（原）
	service.Init()

	_ = proto.RegisterDBProxyServiceHandler(service.Server(), new(rpc.DBProxy))
	if err := service.Run(); err != nil {
		logrus.Error(err)
	}
}

func main() {
	xlog.Init()
	startDbRpcService()
}
