package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"xcloud/common"
	cfg "xcloud/config"
	"xcloud/mq"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
	"xcloud/store/ceph"
	"xcloud/store/oss"
	"xcloud/utils"
)

// UploadGetHandler: 返回上传页面
func UploadGetHandler(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")
	host := c.Request.Host
	ip := strings.Split(host, ":")[0]
	c.HTML(http.StatusOK, "upload.html", gin.H{
		"username": username,
		"token": token,
		"ip": ip,
	})
}

// UploadHandler: 处理文件上传
func UploadPostHandler(c *gin.Context) {
	resp := utils.Resp{
		Code: 0,
		Msg:  "",
	}
	defer func() {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.JSON(http.StatusOK, resp)
	}()
	// 1.从form表单中获得文件内容句柄
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logrus.Error(err)
		resp.Code = common.StatusReadFileError
		resp.Msg = "读取上传的文件出错"
		return
	}
	defer file.Close()
	// 2.把文件内容转换为 []byte
	bufFile := bytes.NewBuffer(nil)
	if _, err := io.Copy(bufFile, file); err != nil {
		logrus.Error(err)
		resp.Code = common.StatusError
		resp.Msg = "转换文件出错"
		return
	}
	// 3.构建文件元信息
	var tempLocalRootDir = cfg.Viper.GetString("file.TempLocalRootDir")
	fileHash := utils.Sha1(bufFile.Bytes())
	fileName := header.Filename
	fileSize := int64(len(bufFile.Bytes()))
	fileAddr := tempLocalRootDir + fileHash // 临时存储地址
	fileMeta := mapper.FileMeta{
		FileHash: fileHash,
		FileName: fileName,
		FileSize: fileSize,
		FileAddr: fileAddr,
	}
	// 4.将文件写入临时存储位置
	tempFile, err := os.Create(fileAddr)
	if err != nil {
		logrus.Errorf("failed to create file, err: %s", err)
		resp.Code = common.StatusServerError
		resp.Msg = "创建文件出错"
		return
	}
	defer tempFile.Close()
	byteCount, err := tempFile.Write(bufFile.Bytes())
	if int64(byteCount) != fileSize || err != nil {
		logrus.Errorf("failed to save data to file; written size: %d, err: %s", byteCount, err)
		resp.Code = common.StatusServerError
		resp.Msg = "写入文件出错"
		return
	}
	// 5.同步/异步地将文件转移到ceph/oss
	_, _ = tempFile.Seek(0, 0) // 文件游标从新回到文件头部
	var currStoreType = cfg.Viper.GetInt("file.CurrStoreType")
	if currStoreType == int(common.StoreCeph) {
		// 文件存储到ceph
		var cephRootDir = cfg.Viper.GetString("ceph.CephRootDir")
		var cephBucket = cfg.Viper.GetString("ceph.CephBucket")
		data, _ := ioutil.ReadAll(tempFile)
		cephFileAddr := cephRootDir + fileHash
		err = ceph.PutObject(cephBucket, cephFileAddr, data)
		if err != nil {
			logrus.Error(err)
			resp.Code = common.StatusError
			resp.Msg = "转移文件至ceph出错"
			return
		}
		// 成功，更新文件元信息中的存储位置
		fileMeta.FileAddr = cephFileAddr
	} else if currStoreType == int(common.StoreOSS) {
		// 文件存储到阿里云oss
		var ossRootDir = cfg.Viper.GetString("oss.OSSRootDir")
		var ossFileAddr string
		strs := strings.Split(fileName, ".")
		if len(strs) >= 2 {
			suffix := strs[len(strs)-1]
			ossFileAddr = ossRootDir + fileHash + "." + suffix
		} else {
			ossFileAddr = ossRootDir + fileHash
		}
		// 判断写入oss为同步还是异步
		var async = cfg.Viper.GetBool("rabbitmq.AsyncTransferEnable")
		if !async {
			// 同步写入
			// TODO: 设置oss中的文件名，方便指定文件名下载
			var ossBucket = cfg.Viper.GetString("oss.Bucket")
			err = oss.GetBucket(ossBucket).PutObject(ossRootDir, tempFile)
			if err != nil {
				logrus.Error(err)
				resp.Code = common.StatusError
				resp.Msg = "转移文件至阿里云oss出错"
				return
			}
			// 成功，更新文件元信息中的存储位置
			fileMeta.FileAddr = ossFileAddr
		} else {
			logrus.Info("文件转移消息发送至rabbitmq中...")
			// 异步写入，写入异步转移任务队列
			transferData := mq.TransferData{
				FileHash:     fileHash,
				CurrLocation: fileAddr,
				DstLocation:  ossFileAddr,
				DstStoreType: common.StoreOSS,
			}
			publishData, _ := json.Marshal(transferData)
			var transExchangeName = cfg.Viper.GetString("rabbitmq.TransExchangeName")
			var transOSSRoutingKey = cfg.Viper.GetString("rabbitmq.TransOSSRoutingKey")
			success := mq.Publish(transExchangeName, transOSSRoutingKey, publishData)
			if !success {
				logrus.Error("当前发送转移消息失败")
				// TODO: 当前发送转移消息失败，稍后重试
			}
		}
	}
	// 6.更新唯一文件表记录
	sqlResult := dbproxy.FileUploadFinished(fileMeta)
	if !sqlResult.Succ {
		resp.Code = common.StatusError
		resp.Msg = sqlResult.Msg
		return
	}
	// 7.更新用户文件表记录
	username := c.Query("username")
	sqlResult = dbproxy.UserFileUploadFinished(username, fileMeta)
	if !sqlResult.Succ {
		resp.Code = common.StatusError
		resp.Msg = sqlResult.Msg
		return
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	return
}

// FastUploadGetHandler: 返回上传页面
func FastUploadGetHandler(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")
	host := c.Request.Host
	ip := strings.Split(host, ":")[0]
	c.HTML(http.StatusOK, "fastupload.html", gin.H{
		"username": username,
		"token": token,
		"ip": ip,
	})
}

// FastUploadPostHandler: 秒传接口
func FastUploadPostHandler(c *gin.Context) {
	username := c.Query("username")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	// 查询相同hash的文件记录
	sqlResult := dbproxy.GetFileMeta(fileHash)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusError,
			"msg":  "秒传失败，请访问普通上传接口",
		})
		return
	}
	fileInfo := mapper.ToFileInfo(sqlResult.Data)
	fileSize := fileInfo.FileSize
	fileAddr := fileInfo.FileAddr
	fileMeta := mapper.FileMeta{
		FileHash: fileHash,
		FileName: fileName,
		FileSize: fileSize,
		FileAddr: fileAddr,
	}

	// 上传过则将文件信息写入用户文件表，返回成功
	sqlResult = dbproxy.UserFileUploadFinished(username, fileMeta)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusError,
			"msg":  "秒传失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusOK,
		"msg":  "秒传成功",
	})
	return
}
