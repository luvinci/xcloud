package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	cfg "xcloud/config"
)

var (
	Endpoint        = cfg.Viper.GetString("oss.Endpoint")
	AccesskeyID     = cfg.Viper.GetString("oss.AccesskeyID")
	AccessKeySecret = cfg.Viper.GetString("oss.AccessKeySecret")
)

var ossCli *oss.Client

// Client: 创建oss client对象
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(Endpoint, AccesskeyID, AccessKeySecret)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return ossCli
}

// GetBucket: 获取bucket存储空间
func GetBucket(bucket string) *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(bucket)
		if err != nil {
			logrus.Error(err)
			return nil
		}
		return bucket
	}
	return nil
}

// DownloadUrl: 生成临时授权下载url
func DownloadUrl(bucket string, objName string) string {
	signUrl, err := GetBucket(bucket).SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return signUrl
}

// 指定为oss的某个bucket设定过期规则
func SetBucketLifecycle(bucketName string) {
	// 表示前缀为test的对象(文件)距最后修改时间30天后过期
	rule1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{rule1}
	err := Client().SetBucketLifecycle(bucketName, rules)
	if err != nil {
		logrus.Infof("Set rule for bucket: %v failed, %v", bucketName, err)
	}
}
