package route

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"xcloud/service/apigw/auth"
	"xcloud/service/apigw/handler"
	"xcloud/utils"
)


func Router() *gin.Engine {
	r := gin.Default()
	// 将静态文件打包成bin文件
	r.Use(static.Serve("/static/", utils.BinaryFileSystem("static")))
	r.LoadHTMLFiles("./static/view/index.html")
	// r.Static("/static/", "./static")

	r.GET("/user/signup", handler.SignupGetHandler)
	r.POST("/user/signup", handler.SignupPostHandler)
	r.GET("/user/signin", handler.SigninGetHandler)
	r.POST("/user/signin", handler.SigninPostHandler)

	r.Use(auth.Authorize())

	r.GET("/", handler.IndexHandler)
	// 获取用户信息
	r.GET("/user/info", handler.GetUserInfoHandler)
	// 获取用户文件列表
	r.POST("/user/files", handler.GetUserFileListHandler)
	// 修改用户文件名
	r.POST("/file/rename", handler.RenameFileHandler)
	// 删除用户文件
	r.POST("/file/delete", handler.DeleteFileHandler)

	return r
}
