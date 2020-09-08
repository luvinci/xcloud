package rpc

import (
	"context"
	"xcloud/common"
	"xcloud/service/account/proto"
	dbproxy "xcloud/service/dbproxy/client"
)

func (*User) DeleteFile(ctx context.Context, req *proto.DeleteFileReq, resp *proto.DeleteFileResp) error {
	fileHash := req.Filehash
	sqlResult := dbproxy.DeleteFile(fileHash)
	if !sqlResult.Succ {
		resp.Code = common.StatusServerError
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	return nil
}
