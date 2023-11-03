package config

import (
	"LDFS/nodeClient"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	ShardsDir      string
	NameNodeClient *nodeClient.NameNodeHttpClient
	DataNodeClient *nodeClient.DataNodeHttpClient

	AddrHttpLocalServer string

	RecoveringShardHash *sync.Map
)

//初始化配置读取器viper
func ConfigInit() error {
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
	AddrHttpLocalServer = viper.GetString("node.http_ip")

	RecoveringShardHash = new(sync.Map)
	ShardsDir = viper.GetString("Data.shardsDir")

	// 种子随机数生成器
	rand.Seed(time.Now().UnixNano())

	return err
}
