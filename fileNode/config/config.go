package config

import (
	"LDFS/nodeClient"
	"fmt"

	"github.com/spf13/viper"
)

var (
	ShardsDir      string
	NameNodeClient *nodeClient.NameNodeHttpClient
	DataNodeClient *nodeClient.DataNodeHttpClient
)

//初始化配置读取器viper
func ViperInit() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadConfig() failed : %s", err)
		return err
	}

	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()

	ShardsDir = viper.GetString("Data.shardsDir")
	return err
}
