package main

import (
	xlog "github.com/luvinci/x-logrus"
	cfg "xcloud/config"
	"xcloud/service/apigw/route"
)

var port = cfg.Viper.GetString("app.apigw.port")

func main() {
	xlog.Init()
	e := route.Router()
	err := e.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
