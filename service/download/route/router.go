package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"xcloud/service/download/handler"
	"xcloud/utils"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(static.Serve("/static/", utils.BinaryFileSystem("static")))
	// r.Static("/static/", "./static")

	// 使用gin插件支持跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"}
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
	}))

	// r.Use(auth.Authorize())

	// 文件上传相关接口
	r.POST("/file/downloadurl", handler.DownloadUrlHandler)
	r.GET("/file/download", handler.DownloadFileHandler)

	return r
}
