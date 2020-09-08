package mapper

import (
	"errors"
	"reflect"
)

var funcs = map[string]interface{} {
	"/user/UserSignUp": UserSignUp,
	"/user/UserSignIn": UserSignIn,
	"/user/UserExist": UserExist,
	"/user/GetPasswordSalt": GetPasswordSalt,
	"/user/UpdateUserToken": UpdateUserToken,
	"/user/SaveTokenToRedis": SaveTokenToRedis,
	"/user/GetUserInfo": GetUserInfo,
	"/user/GetUserFiles": GetUserFiles,

	"/file/RenameUserFile": RenameUserFile,
	"/file/DeleteUserFile": DeleteUserFile,
	"/file/DeleteFile": DeleteFile,
	"/file/DeleteUserFileAndUniqueFile": DeleteUserFileAndUniqueFile,
	"/file/GetSameFileHashCount": GetSameFileHashCount,
	"/file/GetFileAddr": GetFileAddr,
	"/file/GetFileMeta": GetFileMeta,

	"/file/FileUploadFinished": FileUploadFinished,
	"/file/UserFileUploadFinished": UserFileUploadFinished,
	"/file/UpdateFileAddr": UpdateFileAddr,
}

func FuncCall(name string, params ...interface{}) ([]reflect.Value, error) {
	if _, ok := funcs[name]; !ok {
		return nil, errors.New("指定调用的函数不存在")
	}
	// 通过反射可以动态调用对象的导出方法
	f := reflect.ValueOf(funcs[name])
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("传入参数与调用函数的参数长度不一致")
	}
	// 构造一个Value的slice, 用作Call()方法的传入参数
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	// 执行方法f, 将结果返回
	return f.Call(in), nil
}
