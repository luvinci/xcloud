package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
	"xcloud/common"
	cfg "xcloud/config"
	userProto "xcloud/service/account/proto"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
	"xcloud/store/oss"
	"xcloud/utils"
)

func GetUserFileListHandler(c *gin.Context) {
	username := c.Query("username")
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	resp, _ := usercli.GetUserFiles(context.TODO(), &userProto.UserFilesReq{
		Username: username, Limit: int32(limit)})
	if resp.Code < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	var userFiles []mapper.UserFile
	_ = json.Unmarshal(resp.Data, &userFiles)
	response := utils.Resp{
		Code: int(resp.Code),
		Msg:  resp.Msg,
		Data: userFiles,
	}
	// c.Data(http.StatusOK, "application/json", resp.Data)
	c.JSON(http.StatusOK, response)
}

func RenameFileHandler(c *gin.Context) {
	username := c.Query("username")
	fileHash := c.Request.FormValue("filehash")
	fileName := c.Request.FormValue("filename")
	newFileName := c.Request.FormValue("newfilename")
	strs := strings.Split(fileName, ".")
	var suffix string
	var newname string
	if len(strs) >= 2 {
		suffix = strs[len(strs)-1]
		newname = newFileName + "." + suffix
	} else {
		newname = newFileName
	}
	resp, _ := usercli.RenameUserFile(context.TODO(), &userProto.RenameUserFileReq{
		Username: username, Filehash: fileHash, NewFilename: newname})
	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}

func DeleteFileHandler(c *gin.Context) {
	username := c.Query("username")
	fileHash := c.Request.FormValue("filehash")
	/*
		1.根据filehash查找user_file表中的该文件记录个数
		2.如果只有一条记录，说明只有该用户拥有此文件
			获取该文件存储路径
			删除文件记录
			删除本地实体文件
			删除oss上的文件（根据文件存储路径判断）
		3.如果有多条记录，删除该用户文件记录即可
	*/
	sqlResult := dbproxy.GetSameFileHashCount(fileHash)
	if !sqlResult.Succ {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusError,
			"msg":  sqlResult.Msg,
		})
		return
	}
	count := int(sqlResult.Data.(float64))
	if count == 1 {
		sqlResult = dbproxy.GetFileAddr(fileHash)
		if !sqlResult.Succ {
			c.JSON(http.StatusOK, gin.H{
				"code": common.StatusError,
				"msg":  sqlResult.Msg,
			})
			return
		}
		fileAddr := sqlResult.Data.(string)
		// resp, _ := usercli.DeleteUserFileAndUniqueFile(context.TODO(), &userProto.DeleteAllReq{Username: username, Filehash: fileHash})
		resp, _ := usercli.DeleteUserFile(context.TODO(), &userProto.DeleteUserFileReq{Username: username, Filehash: fileHash})
		if resp.Code < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": resp.Code,
				"msg":  resp.Msg,
			})
			return
		}
		resp1, _ := usercli.DeleteFile(context.TODO(), &userProto.DeleteFileReq{Filehash: fileHash})
		if resp.Code < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": resp1.Code,
				"msg":  resp1.Msg,
			})
			return
		}
		prefix := strings.Split(fileAddr, "/")[0]
		var ossRootDir = cfg.Viper.GetString("oss.OSSRootDir")
		ossPrefix := strings.Split(ossRootDir, "/")[0]
		if prefix == ossPrefix {
			var ossBucket = cfg.Viper.GetString("oss.Bucket")
			_ = oss.GetBucket(ossBucket).DeleteObject(fileAddr)
			var tempLocalRootDir = cfg.Viper.GetString("file.TempLocalRootDir")
			fileLocalAddr := tempLocalRootDir + fileHash
			_ = os.Remove(fileLocalAddr)
		} else {
			_ = os.Remove(fileAddr)
		}
	} else {
		resp, _ := usercli.DeleteUserFile(context.TODO(), &userProto.DeleteUserFileReq{Username: username, Filehash: fileHash})
		if resp.Code < 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": resp.Code,
				"msg":  resp.Msg,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": common.StatusOK,
		"msg":  "删除成功",
	})
}
