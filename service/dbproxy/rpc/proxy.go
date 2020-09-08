package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"xcloud/service/dbproxy/mapper"
	"xcloud/service/dbproxy/proto"
)

type DBProxy struct {}

func (p *DBProxy) ExecAction(ctx context.Context, req *proto.ExecReq, resp *proto.ExecResp) error { // 这里此error永远都是nil，因为err都放在SqlResult.Msg中返回去
	/*
	Action  -->  []*SingleAction

	type SingleAction struct {
		Name   string	// Name为字符串，通过反射即可拿到其对应的要执行访问数据库的方法
		Params []byte	// 访问数据库的方法的参数
	}
	*/

	// sqlResults保存着访问数据库的结果（客户端client.go对应的方法即可拿到这个结果）
	sqlResults := make([]mapper.SqlResult, len(req.Action))
	// TODO 检查 req.Sequence req.Transaction两个参数，执行不同的流程（比如是否开启事务）
	for idx, singleAction := range req.Action {
		params := make([]interface{}, 0)
		decoder := json.NewDecoder(bytes.NewReader(singleAction.Params))
		decoder.UseNumber()
		if err := decoder.Decode(&params); err != nil {
			sqlResults[idx] = mapper.SqlResult{
				Succ: false,
				Msg:  "请求参数有误",
			}
			continue
		}
		// 将json.Number类型的参数统一转为int64，避免出现非必要的float64
		// 注意：如果请求的参数为float类型，这里将会被置成0，所以请求参数，比如分数、金额，这种，请使用string类型表示
		for k, param := range params {
			if v, ok := param.(json.Number); ok {
				params[k], _ = v.Int64()
			}
		}

		SqlResult, err := mapper.FuncCall(singleAction.Name, params...)
		if err != nil {
			sqlResults[idx] = mapper.SqlResult{
				Succ: false,
				Msg:  err.Error(),
			}
			continue
		}
		sqlResults = append(sqlResults, SqlResult[0].Interface().(mapper.SqlResult))
	}
	datas, _ := json.Marshal(sqlResults)
	resp.Msg = "执行成功"
	resp.Data = datas
	return nil
}
