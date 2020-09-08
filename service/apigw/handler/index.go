package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func IndexHandler(c *gin.Context) {
	// c.Redirect(http.StatusFound, "/static/view/index.html")
	host := c.Request.Host
	ip := strings.Split(host, ":")[0]
	c.HTML(http.StatusOK, "index.html", gin.H{
		"ip": ip,
	})
}
