package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	cfg "xcloud/config"
)

var conn *amqp.Connection
var channel *amqp.Channel
var notifyClose chan *amqp.Error // 如果异常关闭，会接收通知

var (
	AsyncTransferEnable = cfg.Viper.GetBool("rabbitmq.AsyncTransferEnable")
	URL = cfg.Viper.GetString("rabbitmq.URL")
)

func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitmq连接
	if !AsyncTransferEnable {
		return
	}
	if initChannel(URL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				logrus.Warnf("Notify channel close: %v", msg)
				initChannel(URL)
			}
		}
	}()
}

func initChannel(url string) bool {
	// 1.判断 channel 是否已经创建过
	if channel != nil {
		return true
	}
	// 2.获取 rabbitmq 的一个连接
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		logrus.Error(err)
		return false
	}
	// 3.打开一个 channel 用于消息的发布与接收
	channel, err = conn.Channel()
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}