package client

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	cfg "xcloud/config"

	// cfg "xcloud/config"
	"xcloud/service/dbproxy/mapper"
	"xcloud/service/dbproxy/proto"
)

// type FileMeta struct {
// 	FileHash string
// 	FileName string
// 	FileSize int64
// 	Location string
// 	UploadAt string
// }

var (
	dbcli proto.DBProxyService
	consulAddr = cfg.Viper.GetString("consul.addr")
)

func init() {
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
	)
	// 初始化，解析命令行参数等
	service.Init()
	// 初始化一个dbproxy服务的客户端
	dbcli = proto.NewDBProxyService("go.micro.service.dbproxy", service.Client())
}

// execAction: 向dbproxy请求执行action
func execAction(funcName string, params []byte) *proto.ExecResp {
	// 调用了rpc/proxy.go中的ExecAction方法
	execResp, _ :=  dbcli.ExecAction(
		context.TODO(),
		&proto.ExecReq{
			Action: []*proto.SingleAction{
				{Name: funcName, Params: params},
			},
		})
	return execResp
}

// parseExecRespToSqlResult: 转换rpc返回的结果
func parseExecRespToSqlResult(resp *proto.ExecResp) *mapper.SqlResult {
	if resp == nil || resp.Data == nil {
		return nil
	}
	results := make([]mapper.SqlResult, 0)
	err := json.Unmarshal(resp.Data, &results)
	if err != nil {
		return nil
	}
	if len(results) > 0 {
		return &results[1]
	}
	return nil
}

func UserSignUp(username string, password string, encPassword string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, password, encPassword})
	execResp := execAction("/user/UserSignUp", params)
	return parseExecRespToSqlResult(execResp)
}

func UserSignIn(username string, password string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, password})
	execResp := execAction("/user/UserSignIn", params)
	return parseExecRespToSqlResult(execResp)
}

func UserExist(username string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username})
	execResp := execAction("/user/UserExist", params)
	return parseExecRespToSqlResult(execResp)
}

func GetPasswordSalt(username string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username})
	execResp := execAction("/user/GetPasswordSalt", params)
	return parseExecRespToSqlResult(execResp)
}

func UpdateUserToken(username string, token string) *mapper.SqlResult{
	params, _ := json.Marshal([]interface{}{username, token})
	execResp := execAction("/user/UpdateUserToken", params)
	return parseExecRespToSqlResult(execResp)
}

func SaveTokenToRedis(op, username, expireTime, token string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{op, username, expireTime, token})
	execResp := execAction("/user/SaveTokenToRedis", params)
	return parseExecRespToSqlResult(execResp)
}

func GetUserInfo(username string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username})
	execResp := execAction("/user/GetUserInfo", params)
	return parseExecRespToSqlResult(execResp)
}

func GetUserFiles(username string, limit int) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, limit})
	execResp := execAction("/user/GetUserFiles", params)
	return parseExecRespToSqlResult(execResp)
}

func RenameUserFile(username string, fileHash string, newFilename string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, fileHash, newFilename})
	execResp := execAction("/file/RenameUserFile", params)
	return parseExecRespToSqlResult(execResp)
}

func DeleteUserFile(username string, fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, fileHash})
	execResp := execAction("/file/DeleteUserFile", params)
	return parseExecRespToSqlResult(execResp)
}

func GetSameFileHashCount(fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileHash})
	execResp := execAction("/file/GetSameFileHashCount", params)
	return parseExecRespToSqlResult(execResp)
}

func DeleteFile(fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileHash})
	execResp := execAction("/file/DeleteFile", params)
	return parseExecRespToSqlResult(execResp)
}

func DeleteUserFileAndUniqueFile(username string, fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, fileHash})
	execResp := execAction("/file/DeleteUserFileAndUniqueFile", params)
	return parseExecRespToSqlResult(execResp)
}

func GetFileAddr(fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileHash})
	execResp := execAction("/file/GetFileAddr", params)
	return parseExecRespToSqlResult(execResp)
}

func GetFileMeta(fileHash string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileHash})
	execResp := execAction("/file/GetFileMeta", params)
	return parseExecRespToSqlResult(execResp)
}

// FileUploadFinished: 文件上传成功，新增/更新（唯一文件）元信息
func FileUploadFinished(fileMeta mapper.FileMeta) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileMeta.FileHash, fileMeta.FileName, fileMeta.FileAddr, fileMeta.FileSize})
	execResp := execAction("/file/FileUploadFinished", params)
	return parseExecRespToSqlResult(execResp)
}

// UserFileUploadFinished: 文件上传成功，新增/更新（用户文件）元信息
func UserFileUploadFinished(username string, fileMeta mapper.FileMeta) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{username, fileMeta.FileHash, fileMeta.FileName, fileMeta.FileSize})
	execResp := execAction("/file/UserFileUploadFinished", params)
	return parseExecRespToSqlResult(execResp)
}

// UpdateFileAddr: 更新文件存储路径（文件转移后）
func UpdateFileAddr(fileHash string, fileAddr string) *mapper.SqlResult {
	params, _ := json.Marshal([]interface{}{fileHash, fileAddr})
	execResp := execAction("/file/UpdateFileAddr", params)
	return parseExecRespToSqlResult(execResp)
}




























