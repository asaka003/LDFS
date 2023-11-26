package config

import (
	"LDFS/nodeClient"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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
	ShardsDir = filepath.Join(viper.GetString("Data.shardsDir"), viper.GetString("node.name"))
	//创建目录
	_, err = os.Stat(ShardsDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(ShardsDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 种子随机数生成器
	rand.Seed(time.Now().UnixNano())

	return err
}
