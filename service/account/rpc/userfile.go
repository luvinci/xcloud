package rpc

import (
	"context"
	"encoding/json"
	"xcloud/common"
	"xcloud/service/account/proto"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
)

func (*User) GetUserFiles(ctx context.Context, req *proto.UserFilesReq, resp *proto.UserFilesResp) error {
	username := req.Username
	limit := req.Limit
	sqlResult := dbproxy.GetUserFiles(username, int(limit))
	if !sqlResult.Succ {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}

	userFiles := mapper.ToUserFile(sqlResult.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	resp.Data = data
	return nil
}

func (*User) RenameUserFile(ctx context.Context, req *proto.RenameUserFileReq, resp *proto.RenameUserFileResp) error {
	username := req.Username
	fileHash := req.Filehash
	newFilename := req.NewFilename
	sqlResult := dbproxy.RenameUserFile(username, fileHash, newFilename)
	if !sqlResult.Succ {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	return nil
}

func (*User) DeleteUserFile(ctx context.Context, req *proto.DeleteUserFileReq, resp *proto.DeleteUserFileResp) error {
	username := req.Username
	fileHash := req.Filehash
	sqlResult := dbproxy.DeleteUserFile(username, fileHash)
	if !sqlResult.Succ {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	return nil
}

func (*User)DeleteUserFileAndUniqueFile(ctx context.Context, req *proto.DeleteAllReq, resp *proto.DeleteAllResp) error {
	username := req.Username
	fileHash := req.Filehash
	sqlResult := dbproxy.DeleteUserFile(username, fileHash)
	if !sqlResult.Succ {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	return nil
}