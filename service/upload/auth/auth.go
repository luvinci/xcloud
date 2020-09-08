package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"net/http"
	rds "xcloud/cache/redis"
)

// Authorize: http请求拦截验证token
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		token := c.Query("token")
		if err := IsTokenValid(username, token); err != nil  {
			logrus.Error(err)
			// w.WriteHeader(http.StatusForbidden)
			// token校验失败则跳转到登录页面
			c.Redirect(http.StatusFound, "http://localhost:8080/user/signin")
			c.Abort()
		}
		c.Next()
	}
}

// IsTokenValid: token是否有效
func IsTokenValid(username string, token string) error {
	conn := rds.Pool().Get()
	defer conn.Close()
	res, _ := redis.String(conn.Do("get", username))
	if res != token {
		return errors.New("token校验失败")
	}
	// token校验成功，重置token过期时间
	_, _ = conn.Do("setex", username, "1200", token)
	return nil
}
