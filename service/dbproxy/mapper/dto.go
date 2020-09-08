package mapper

import (
	"github.com/mitchellh/mapstructure"
)

// SqlResult: 执行sql返回的结果
type SqlResult struct {
	Succ bool        `json:"succ"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type UserInfo struct {
	UserName     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

type UserFile struct {
	FileHash   string
	FileName   string
	FileSize   int64
	UploadAt   string
	LastUpdate string
}

type FileInfo struct {
	FileSize int64
	FileAddr string
}

func ToUserInfo(i interface{}) UserInfo {
	userInfo := UserInfo{}
	_ = mapstructure.Decode(i, &userInfo)
	return userInfo
}

func ToUserFile(i interface{}) []UserFile {
	var userFiles []UserFile
	_ = mapstructure.Decode(i, &userFiles)
	return userFiles
}

func ToFileInfo(i interface{}) FileInfo {
	fileInfo := FileInfo{}
	_ = mapstructure.Decode(i, &fileInfo)
	return fileInfo
}
