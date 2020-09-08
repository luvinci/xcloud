package mq

import (
	"xcloud/common"
)

// TransferData: rabbitmq转移队列中消息载体的结构
type TransferData struct {
	FileHash     string
	CurrLocation string           // 文件临时存储的地址
	DstLocation  string           // 要转移的目标地址
	DstStoreType common.StoreType // 文件存储类型
}
