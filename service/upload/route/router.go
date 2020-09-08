package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"xcloud/service/upload/auth"
	"xcloud/service/upload/handler"
	"xcloud/utils"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(static.Serve("/static/", utils.BinaryFileSystem("static")))
	r.LoadHTMLFiles("./static/view/upload.html", "./static/view/fastupload.html")
	// r.Static("/static/", "./static")

	// 使用gin插件支持跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"}, // []string{"http://localhost:8080"}
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
	}))

	r.Use(auth.Authorize())

	// 文件上传相关接口
	r.GET("/file/upload", handler.UploadGetHandler)
	r.POST("/file/upload", handler.UploadPostHandler)

	// 秒传接口
	r.GET("/file/fastupload", handler.FastUploadGetHandler)
	r.POST("/file/fastupload", handler.FastUploadPostHandler)

	// 分块上传接口
	r.POST("/file/mpupload/upinit", handler.UploadInitHandler)
	r.POST("/file/mpupload/upchunk", handler.UploadChunkHandler)
	r.POST("/file/mpupload/upcomplete", handler.UploadCompleteHandler)

	return r
}
