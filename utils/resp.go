package utils

import (
	"encoding/json"
)

// Resp: http响应数据的通用结构
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// ToBytes: 对象转json格式的二进制数组
func (resp *Resp) ToJsonBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		return nil
	}
	return r
}

// ToString: 对象转json格式的string
func (resp *Resp) ToJsonString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		return ""
	}
	return string(r)
}