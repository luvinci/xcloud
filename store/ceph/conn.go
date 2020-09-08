package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"
	cfg "xcloud/config"
)

var (
	AccessKey   = cfg.Viper.GetString("ceph.AccessKey")
	SecretKey   = cfg.Viper.GetString("ceph.SecretKey")
	EC2Endpoint = cfg.Viper.GetString("ceph.EC2Endpoint")
	S3Endpoint  = cfg.Viper.GetString("ceph.S3Endpoint")
)

var cephConn *s3.S3

// Conn: 获取ceph连接
func Conn() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}
	// 1.初始化ceph的一些信息
	auth := aws.Auth{
		AccessKey: AccessKey,
		SecretKey: SecretKey,
	}
	region := aws.Region{
		Name:                 "default",
		EC2Endpoint:          EC2Endpoint,
		S3Endpoint:           S3Endpoint,
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}
	// 2.创建s3类型的连接
	return s3.New(auth, region)
}

// GetCephBucket: 获取指定的bucket对象
func GetCephBucket(bucket string) *s3.Bucket {
	conn := Conn()
	return conn.Bucket(bucket)
}

// PutObject: 上传文件到ceph集群
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}