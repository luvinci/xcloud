package rpc

import (
	"context"
	"xcloud/common"
	"xcloud/service/account/proto"
	dbproxy "xcloud/service/dbproxy/client"
	"xcloud/service/dbproxy/mapper"
	"xcloud/utils"
)

type User struct {}

func (*User) SignUp(ctx context.Context, req *proto.SignUpReq, resp *proto.SignUpResp) error {
	username := req.Username
	password := req.Password
	// 判断username是否已经注册
	sqlResult := dbproxy.UserExist(username)
	if sqlResult.Succ {
		resp.Code = common.StatusUserExists
		resp.Msg = sqlResult.Msg
		return nil
	}

	// 生成密码盐值，对密码进行加密
	passwordSalt := utils.GeneratePasswordSalt()
	encPassword := utils.EncryptPassword(password, passwordSalt)
	sqlResult = dbproxy.UserSignUp(username, encPassword, passwordSalt)
	if !sqlResult.Succ {
		resp.Code = common.StatusRegisterFailed
		resp.Msg = sqlResult.Msg
		return nil
	}
	resp.Code = common.StatusRegisterSuccess
	resp.Msg = "注册成功"
	return nil
}

func (*User) SignIn(ctx context.Context, req *proto.SignInReq, resp *proto.SignInResp) error {
	username := req.Username
	password := req.Password
	// 1.从数据获取密码盐值，校验用户名和密码
	sqlResult := dbproxy.GetPasswordSalt(username)
	if !sqlResult.Succ {
		resp.Code = common.StatusUserNotExists
		resp.Msg = sqlResult.Msg
		return nil
	}
	passwordSalt := sqlResult.Data.(string)
	encPassword := utils.EncryptPassword(password, passwordSalt)
	sqlResult = dbproxy.UserSignIn(username, encPassword)
	if !sqlResult.Succ {
		resp.Code = common.StatusError
		resp.Msg = sqlResult.Msg
		return nil
	}
	// 2.sqlResult.Succ为true，通过认证，生成访问凭证（token）保存
	token := utils.GenerateToken(username)
	sqlResult = dbproxy.UpdateUserToken(username, token)
	if !sqlResult.Succ {
		resp.Code = common.StatusError
		resp.Msg = sqlResult.Msg
		return nil
	}
	// 3.保存token到redis，设置过期时间
	sqlResult = dbproxy.SaveTokenToRedis("setex", username, "1200", token)
	if !sqlResult.Succ {
		resp.Code = common.StatusSetTokenFailed
		resp.Msg = sqlResult.Msg
		return nil
	}
	// 4.返回token
	resp.Code = common.StatusLoginSuccess
	resp.Msg = sqlResult.Msg
	resp.Token = token
	return nil
}

func (*User) GetUserInfo(ctx context.Context, req *proto.UserInfoReq, resp *proto.UserInfoResp) error {
	username := req.Username
	sqlResult := dbproxy.GetUserInfo(username)
	if !sqlResult.Succ {
		resp.Code = common.StatusError
		resp.Msg = sqlResult.Msg
		return nil
	}
	// 组装并且响应用户数据
	userInfo := mapper.ToUserInfo(sqlResult.Data)
	resp.Code = common.StatusOK
	resp.Msg = sqlResult.Msg
	resp.Username = userInfo.UserName
	// TODO: 需增加接口支持完善用户信息(email/phone等)
	resp.Email = userInfo.Email
	resp.Phone = userInfo.Phone
	resp.SignupAt = userInfo.SignupAt
	resp.LastActiveAt = userInfo.LastActiveAt
	resp.Status = int32(userInfo.Status)
	return nil
}
