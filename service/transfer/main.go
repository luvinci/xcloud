package main

import (
	xlog "github.com/luvinci/x-logrus"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/sirupsen/logrus"
	"time"
	cfg "xcloud/config"
	"xcloud/mq"
	"xcloud/service/transfer/process"
)

var (
	consulAddr          = cfg.Viper.GetString("consul.addr")
	asyncTransferEnable = cfg.Viper.GetBool("rabbitmq.AsyncTransferEnable")
	transOSSQueueName   = cfg.Viper.GetString("rabbitmq.TransOSSQueueName")
)

func startRpcService() {
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("go.micro.service.transfer"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		// micro.Flags(common.CustomFlags...),
	)
	service.Init()

	if err := service.Run(); err != nil {
		logrus.Error(err)
	}
}

func startTransferservice() {
	if !asyncTransferEnable {
		logrus.Error("异步转移文件功能被禁用，请检查相关设置")
		return
	}
	logrus.Info("文件转移服务启动中，开始监听转移任务队列...")
	mq.StartConsume(transOSSQueueName, "transfer_oss", process.Transfer)
}

func main() {
	xlog.Init()
	mq.Init()
	// 文件转移服务
	go startTransferservice()
	// rpc 服务
	startRpcService()
}
