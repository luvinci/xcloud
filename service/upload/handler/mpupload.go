package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
	rds "xcloud/cache/redis"
	cfg "xcloud/config"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
)

var (
	chunkSize = 2 << 20 // 2MB
	chunkDir  = cfg.Viper.GetString("file.TempChunkRootDir")
	localDir  = cfg.Viper.GetString("file.TempLocalRootDir")
)

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	UploadID   string // 本地存储路径 + UploadID + 分块编号
	FileHash   string
	FileSize   int
	ChunkSize  int // 分块大小
	ChunkCount int // 分块数量
}

// UploadInitHandler: 分块上传初始化工作
func UploadInitHandler(c *gin.Context) {
	username := c.Query("username")
	fileHash := c.Request.FormValue("filehash")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不合格",
		})
		return
	}

	conn := rds.Pool().Get()
	defer conn.Close()

	// 构造初始化信息
	uploadId := username + "_" + fmt.Sprintf("%x", time.Now().UnixNano())
	chunkCount := int(math.Ceil(float64(fileSize)/float64(chunkSize)))
	uploadInfo := MultipartUploadInfo{
		UploadID:   uploadId,
		FileHash:   fileHash,
		FileSize:   fileSize,
		ChunkSize:  chunkSize,
		ChunkCount: chunkCount,
	}
	expireTime := 60 * 60 * 12 // 12小时
	conn.Do("HMSET", "mp_"+uploadId, "filehash", fileHash, "filesize", fileSize, "chunkcount", chunkCount)
	conn.Do("EXPIRE", "mp_"+uploadId, expireTime)
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"data": uploadInfo,
	})
}

// UploadChunkHandler: 分块上传接口
func UploadChunkHandler(c *gin.Context) {
	uploadId := c.Request.FormValue("uploadid")
	chunkIdx := c.Request.FormValue("chunkidx")

	// 创建相应的(分块)文件句柄，用来存储此次客户端上传的分块内容
	chunkIdxFile := chunkDir + uploadId + "/" + chunkIdx
	_ = os.MkdirAll(path.Dir(chunkIdxFile), os.ModePerm)
	chunkFile, err := os.Create(chunkIdxFile)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分块文件失败",
		})
		return
	}
	defer chunkFile.Close()

	buf := make([]byte, chunkSize) // 每次读取2MB
	for {
		n, err := c.Request.Body.Read(buf)
		chunkFile.Write(buf[:n])
		if err != nil {
			break
		}
	}
	defer c.Request.Body.Close()

	conn := rds.Pool().Get()
	defer conn.Close()

	// 更新缓存中此次分块文件所对应的的分块信息
	conn.Do("HSET", "mp_"+uploadId, "chunkidx_"+chunkIdx, 1)


	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "上传分块成功",
	})
}

// UploadCompleteHandler: 分块上传完成，合并分块
func UploadCompleteHandler(c *gin.Context) {
	username := c.Query("username")
	uploadId := c.Request.FormValue("uploadid")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	fileSize, err := strconv.Atoi(c.Request.FormValue("filesize"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不合格",
		})
		return
	}

	conn := rds.Pool().Get()
	defer conn.Close()

	// 通过uploadid查询redis并判断是否所有分块上传完成
	chunkCount := 0
	actualCount := 0
	datas, _ := redis.Values(conn.Do("HGETALL", "mp_"+uploadId))
	for i := 0; i < len(datas); i += 2 {
		k := string(datas[i].([]byte))
		v := string(datas[i+1].([]byte))
		if k == "chunkcount" {
			chunkCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chunkidx_") && v == "1" {
			actualCount++
		}
	}

	if chunkCount != actualCount {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分块不完整",
		})
		return
	}
	// 合并分块到本地
	fileAddr := localDir + fileName
	completeFile, _ := os.Create(fileAddr)

	var chunkFile *os.File
	for i := 1; i <= chunkCount; i++ {
		chunkIdxFile := chunkDir + uploadId + "/" + strconv.Itoa(i)
		chunkFile, err = os.OpenFile(chunkIdxFile, os.O_RDONLY, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "发生未知错误，请稍后重试",
			})
			return
		}
		bytesData, _ := ioutil.ReadAll(chunkFile)
		_, err = completeFile.Write(bytesData)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "发生未知错误，请稍后重试",
			})
			return
		}
		os.Remove(chunkIdxFile)
	}
	completeFile.Close()
	chunkFile.Close()

	// err = os.RemoveAll(chunkDir + uploadId)
	// if err != nil {
	// 	logrus.Error(err)
	// }

	if runtime.GOOS == "linux" {
		err = os.RemoveAll(chunkDir + uploadId)
		if err != nil {
			logrus.Error(err)
		}
	}

	// 更新唯一文件表及用户文件表
	fileMeta := mapper.FileMeta{
		FileHash: fileHash,
		FileName: fileName,
		FileSize: int64(fileSize),
		FileAddr: fileAddr,
	}
	sqlResult := dbproxy.FileUploadFinished(fileMeta)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生未知错误，请稍后重试",
		})
		return
	}
	sqlResult = dbproxy.UserFileUploadFinished(username, fileMeta)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "发生未知错误，请稍后重试",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "上传成功",
	})
	return
}
