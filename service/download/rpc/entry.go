package rpc

import (
	"context"
	cfg "xcloud/config"
	"xcloud/service/download/proto"
)

var DownloadEntry = cfg.Viper.GetString("app.download.entry")

type Download struct {}

func (*Download) DownloadEntry(ctx context.Context, req *proto.EntryReq, resp *proto.EntryResp) error {
	resp.Entry = DownloadEntry
	return nil
}
