package initialize

import (
	"exporter/global"
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {

	v := viper.New()
	v.SetConfigFile("config/config.yaml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatal("配置读取失败:", err)
	}
	if err := v.Unmarshal(&global.GvaServerConfig); err != nil {
		log.Fatal("配置解析失败:", err)
	}
}
