package common

const (
	// 正常
	StatusOK = 10000
	// 注册成功
	StatusRegisterSuccess = 10001
	// 登陆成功
	StatusLoginSuccess = 10002

	// 请求参数无效
	StatusParamInvalid = -10000
	// 服务出错
	StatusServerError = -10001
	// 注册失败
	StatusRegisterFailed = -10002
	// 登录失败
	StatusLoginFailed = -10003
	// token无效
	StatusTokenInvalid = -10004
	// 用户不存在
	StatusUserNotExists = -10005
	// 通用错误
	StatusError = -10006
	// 用户已存在
	StatusUserExists = -10007
	// redis设置token失败
	StatusSetTokenFailed = -10008
	// 读取文件出错
	StatusReadFileError = -10009
)
