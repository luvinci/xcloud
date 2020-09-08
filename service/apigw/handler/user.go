package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	zh "github.com/luvinci/validator-zh"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"net/http"
	"xcloud/common"
	cfg "xcloud/config"
	userProto "xcloud/service/account/proto"
	"xcloud/utils"
)

var (
	usercli    userProto.UserService
	consulAddr = cfg.Viper.GetString("consul.addr")
)

func init() {
	// TODO 置请求容量及qps
	reg := consul.NewRegistry(registry.Addrs(consulAddr))
	service := micro.NewService(
		micro.Registry(reg),
	)
	// 初始化，解析命令行参数等
	service.Init()
	// 初始化一个account服务的客户端
	usercli = userProto.NewUserService("go.micro.service.user", service.Client())
}

type User struct {
	Username string `json:"username" validate:"required,gte=2,lte=16" label:"用户名"`
	Password string `json:"password" validate:"required,min=6,max=128" label:"密码"`
}

// SignupGetHandler: 响应注册页面
func SignupGetHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// SignupPostHandler: 响应注册post请求
func SignupPostHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	// 校验参数
	user := User{
		Username: username,
		Password: password,
	}
	errmsgs := zh.Validate(user)
	if errmsgs != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusParamInvalid,
			"msg":  errmsgs,
		})
		return
	}

	resp, _ := usercli.SignUp(context.TODO(), &userProto.SignUpReq{Username: username, Password: password})
	if resp.Code < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
	})
}

// SigninGetHandler: 响应登陆页面
func SigninGetHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// SigninPostHandler: 响应登陆post请求
func SigninPostHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	// 校验参数
	user := User{
		Username: username,
		Password: password,
	}
	errmsgs := zh.Validate(user)
	if errmsgs != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": common.StatusParamInvalid,
			"msg":  errmsgs,
		})
		return
	}

	resp, _ := usercli.SignIn(context.TODO(), &userProto.SignInReq{Username: username, Password: password})
	if resp.Code < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	data := map[string]string{
		"username": username,
		"token": resp.Token,
		"location": "http://" + c.Request.Host,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": resp.Code,
		"msg":  resp.Msg,
		"data": data,
	})
}

func GetUserInfoHandler(c *gin.Context) {
	username := c.Query("username")
	resp, _ := usercli.GetUserInfo(context.TODO(), &userProto.UserInfoReq{Username: username})
	if resp.Code < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": resp.Code,
			"msg":  resp.Msg,
		})
		return
	}
	utlesResp := utils.Resp{
		Code: int(resp.Code),
		Msg:  resp.Msg,
		Data: gin.H{
			"username": resp.Username,
			"email": resp.Email,
			"phone": resp.Phone,
			"signup_at": resp.SignupAt,
			"last_active_at": resp.LastActiveAt,
			"status": int(resp.Status),
		},
	}
	c.Data(http.StatusOK, "application/json", utlesResp.ToJsonBytes())
}
