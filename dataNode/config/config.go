package config

import (
	"LDFS/nodeClient"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	ShardsDir           string = "LDFS/data-node/data-shards/"
	TempDir             string = "LDFS/data-node/temp-shards/"
	NameNodeClient      *nodeClient.NameNodeHttpClient
	DataNodeClient      *nodeClient.DataNodeHttpClient
	RecoveringShardHash *sync.Map
)

//初始化配置
func ConfigInit() error {
	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()

	RecoveringShardHash = new(sync.Map)
	//创建目录
	_, err := os.Stat(ShardsDir)
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
