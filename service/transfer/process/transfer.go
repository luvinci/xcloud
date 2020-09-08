package process

import (
	"bufio"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	cfg "xcloud/config"
	"xcloud/mq"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/store/oss"
)

var (
	ossBucket = cfg.Viper.GetString("oss.Bucket")
)

func Transfer(msg []byte) bool {
	logrus.Infof("收到转移消息：%s", string(msg))
	transferData := mq.TransferData{}
	err := json.Unmarshal(msg, &transferData)
	if err != nil {
		logrus.Error(err)
		return false
	}
	file, err := os.Open(transferData.CurrLocation)
	if err != nil {
		logrus.Error(err)
		return false
	}
	err = oss.GetBucket(ossBucket).PutObject(
		transferData.DstLocation, bufio.NewReader(file))
	if err != nil {
		logrus.Error(err)
		return false
	}
	sqlResult := dbproxy.UpdateFileAddr(transferData.FileHash, transferData.DstLocation)
	if !sqlResult.Succ {
		logrus.Error(sqlResult.Msg)
		return false
	}
	return true
}
