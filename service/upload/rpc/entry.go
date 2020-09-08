package rpc

import (
	"context"
	cfg "xcloud/config"
	"xcloud/service/upload/proto"
)

var UploadEntry = cfg.Viper.GetString("app.upload.entry")

type Upload struct {}

func (*Upload) UploadEntry(ctx context.Context, req *proto.EntryReq, resp *proto.EntryResp) error {
	resp.Entry = UploadEntry
	return nil
}
