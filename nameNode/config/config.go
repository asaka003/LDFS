package config

import (
	"LDFS/nodeClient"
	"math/rand"
	"time"
)

var (
	DataNodeClient *nodeClient.DataNodeHttpClient

	FileMetaDir string = "LDFS/name-node/meta" //存储文件meta信息的目录

	ECDataShardNum   int64 = 3 //RS校验码参数，默认三个数据快和两个校验块
	ECParityShardNum int64 = 2
	CopyReplicasNum  int64 = 3 //默认副本数量为3
)

const (
	RemainSize int64 = 100 * 1024 * 1024 //系统保留空间，不会被存储系统使用
)

//初始化所有配置信息
func ConfigInit() (err error) {
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	// 种子随机数生成器
	rand.Seed(time.Now().UnixNano())
	return
}

// func printStorageInfo() {
// 	for i, dataNode := range DataNodeList {
// 		fmt.Printf("dataNode%v节点地址:%s\n", i, dataNode.URL)
// 		fmt.Println("free:", dataNode.NodeDiskAvailableSize)
// 		fmt.Println("total:", dataNode.NodeDiskSize)
// 		fmt.Println("used:", dataNode.NodeDiskUsedSize)
// 		fmt.Println("fileTotal:", dataNode.NodeFileTotalSize)
// 		fmt.Println("")
// 	}
// }
