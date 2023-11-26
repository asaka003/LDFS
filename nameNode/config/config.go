package config

import (
	"LDFS/model"
	"LDFS/nodeClient"
	"fmt"
	"math/rand"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	//读取文件个各节点的IP地址
	DataNodeList   []model.DataNode
	DataNodeUrls   []string
	DataNodeClient *nodeClient.DataNodeHttpClient

	MultiUploadDir string
	FileMetaDir    string

	ECDataShardNum   int64
	ECParityShardNum int64
	CopyReplicasNum  int64

	HttpApiServerHost string
)

const (
	RemainSize int64 = 100 * 1024 * 1024
)

//初始化配置读取器viper
func viperInit() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadConfig() failed : %s", err)
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		//设置配置文件修改后的动态绑定事件......
		fmt.Println("配置文件改变")
	})
	return nil
}

//初始化所有配置信息
func ConfigInit() (err error) {
	err = viperInit()
	if err != nil {
		return err
	}
	MultiUploadDir = viper.GetString("MultiUploadDir")
	FileMetaDir = viper.GetString("FileMetaDir")
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	DataNodeUrls = viper.GetStringSlice("DataNodes.List")
	for _, url := range DataNodeUrls {
		//请求DataNode 磁盘存储情况
		dataNode, err := DataNodeClient.GetStorageInfo(url)
		if err != nil {
			return err
		}
		DataNodeList = append(DataNodeList, dataNode)
	}
	printStorageInfo()

	ECDataShardNum = viper.GetInt64("EC.dataShards")
	ECParityShardNum = viper.GetInt64("EC.parityShards")
	CopyReplicasNum = viper.GetInt64("Copy.replicasNum")
	HttpApiServerHost = viper.GetString("HttpApiServerHost")

	// 种子随机数生成器
	rand.Seed(time.Now().UnixNano())

	return
}

func printStorageInfo() {
	for i, dataNode := range DataNodeList {
		fmt.Printf("dataNode%v节点地址:%s\n", i, dataNode.URL)
		fmt.Println("free:", dataNode.NodeDiskAvailableSize)
		fmt.Println("total:", dataNode.NodeDiskSize)
		fmt.Println("used:", dataNode.NodeDiskUsedSize)
		fmt.Println("fileTotal:", dataNode.NodeFileTotalSize)
		fmt.Println("")
	}
}
