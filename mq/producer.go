package mq

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Publish: 发布消息
func Publish(exchange string, routingkey string, msg []byte) bool {
	// 1.判断 channel 是否正常
	if !initChannel(URL) {
		return false
	}
	// 2.执行消息发布动作
	err := channel.Publish(
		exchange,
		routingkey,
		false, // 如果没有对应的queue，交换机是否将消息返回给发布者，false表示不返回（消息丢失）
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}

