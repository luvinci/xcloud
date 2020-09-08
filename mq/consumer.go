package mq

import (
	"github.com/sirupsen/logrus"
)

var done chan struct{}

// StartConsume: 接收消息
func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1.通过 channel.Consume 获得消息信道
	msgs, err := channel.Consume(
		qName, // 队列名称
		cName, // 消费者名称
		true,  // 自动应答通知已收到消息
		false, // 非唯一的消费者，其他消费者处理器也可以去竞争这个队列里面的消息任务
		false, // rabbitMq 不支持了，只能设置false
		false, // false 表示会阻塞直到有消息过来
		nil,
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	done = make(chan struct{})

	// 2.循环获取队列的消息，没有消息就会阻塞
	go func() {
		for msg := range msgs {
			// 3.调用 callback 函數来处理获取的消息
			logrus.Infof("接收到的任务消息：%s", msg.Body)
			success := callback(msg.Body)
			if !success {
				// TODO: 如果任务处理失败，加入错误队列，待后续处理重试
			}
		}
	}()
	// 接收done的信号，没有消息过来就会阻一直阻塞
	<-done
	// 收到消息说明不在监听处理消息，则关闭rabbitmq通道
	_ = channel.Close()
}

// StopConsume: 停止监听队列
func StopConsume() {
	done <- struct{}{}
}
