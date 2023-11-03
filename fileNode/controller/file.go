package controller

import (
	"LDFS/fileNode/config"
	"LDFS/model"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//恢复文件块
func RecoverShard(c *gin.Context) {
	params := new(model.RecoverShardParam)
	err := c.ShouldBindJSON(params)
	if err != nil || strings.Contains(params.Block.Shards[params.ShardId].Hash, "..") {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	//检查该shard是否已经在恢复过程中
	_, ok := config.RecoveringShardHash.Load(params.Block.Shards[params.ShardId].Hash)
	if ok {
		ResponseSuc(c, nil)
		return
	} else {
		config.RecoveringShardHash.Store(params.Block.Shards[params.ShardId].Hash, true)
	}
	AvailableNodeUrls := make([]string, 0)
	for _, shard := range params.Block.Shards {
		if shard.ShardID == params.ShardId {
			continue
		}
		AvailableNodeUrls = append(AvailableNodeUrls, shard.NodeURL)
	}

	shardPath := filepath.Join(config.ShardsDir, params.Block.Shards[params.ShardId].Hash)
	file, err := os.Create(shardPath)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	defer file.Close()
	maxRecoverNum := len(AvailableNodeUrls)
	for i := 0; i < maxRecoverNum; i++ {
		index := rand.Intn(len(AvailableNodeUrls))
		selectUrl := AvailableNodeUrls[index]
		err = config.DataNodeClient.ReplicasDownloadShard(params.Block.Shards[params.ShardId].Hash, selectUrl, file)
		if err == nil {
			break
		}
		AvailableNodeUrls = append(AvailableNodeUrls[:index], AvailableNodeUrls[index+1:]...)
		if i == maxRecoverNum-1 {
			ResponseErr(c, CodeFileError)
			return
		}
	}
	config.RecoveringShardHash.Delete(params.Block.Shards[params.ShardId].Hash)
	ResponseSuc(c, nil)
}
