package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"xcloud/common"
	cfg "xcloud/config"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
	"xcloud/store/ceph"
	"xcloud/store/oss"
)

// DownloadUrlHandler: 生成文件的下载地址
func DownloadUrlHandler(c *gin.Context) {
	username := c.Query("username")
	token := c.Query("token")
	c.Query("token")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	// 查询相同hash的文件记录
	sqlResult := dbproxy.GetFileMeta(fileHash)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusError,
			"msg":  "该文件不存在",
		})
		return
	}
	fileInfo := mapper.ToFileInfo(sqlResult.Data)
	fileAddr := fileInfo.FileAddr
	// 判断文件存在本地，还是oss/ceph
	var localDir = cfg.Viper.GetString("file.TempLocalRootDir")
	var cephDir = cfg.Viper.GetString("ceph.CephRootDir")
	var ossDir = cfg.Viper.GetString("oss.OSSRootDir")
	if strings.HasPrefix(fileAddr, localDir) || strings.HasPrefix(fileAddr, cephDir) {
		tempURL := fmt.Sprintf("http://%s/file/download?filename=%s&fileaddr=%s&username=%s&token=%s",
			c.Request.Host, fileName, fileAddr, username, token)
		c.Data(http.StatusOK, "application/octet-stream", []byte(tempURL))
	} else if strings.HasPrefix(fileAddr, ossDir) {
		// oss下载url
		var bucket = cfg.Viper.GetString("oss.Bucket")
		signedURL := oss.DownloadUrl(bucket, fileAddr)
		c.Data(http.StatusOK, "application/octet-stream", []byte(signedURL))
	}
}

// DownloadFileHandler: 文件下载接口
func DownloadFileHandler(c *gin.Context) {
	fileName := c.Request.FormValue("filename")
	fileAddr := c.Request.FormValue("fileaddr")

	var localDir = cfg.Viper.GetString("file.TempLocalRootDir")
	var cephDir = cfg.Viper.GetString("ceph.CephRootDir")
	if strings.HasPrefix(fileAddr, localDir) {
		// 本地文件，直接下载
		c.FileAttachment(fileAddr, fileName)
	} else if strings.HasPrefix(fileAddr, cephDir) {
		// ceph中的文件，通过ceph api先下载
		var bucket = cfg.Viper.GetString("ceph.CephBucket")
		data, _ := ceph.GetCephBucket(bucket).Get(fileAddr)
		c.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
		c.Data(http.StatusOK, "application/octect-stream", data)
	}
}
