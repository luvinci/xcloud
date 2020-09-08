package config

import (
	"bytes"
	"github.com/spf13/viper"
	"xcloud/assets"
)

var Viper *viper.Viper

/*
func init() {
	Viper = viper.New()
	// 设置配置文件的名字
	Viper.SetConfigName("config")
	// 添加配置文件所在的路径
	Viper.AddConfigPath("./config")
	// 设置配置文件类型
	Viper.SetConfigType("ini")
	if err := Viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
 */

func init() {
	Viper = viper.New()
	bytesData, err := assets.Asset("config/config.ini")
	if err != nil {
		panic(err)
	}
	Viper.SetConfigType("ini")
	if err = Viper.ReadConfig(bytes.NewBuffer(bytesData)); err != nil {
		panic(err)
	}
}