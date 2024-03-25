package controller

import (
	"LDFS/dataNode/config"
	"LDFS/logger"
	"LDFS/model"
	"LDFS/nodeClient"
	storagesdk "LDFS/storage-sdk"
	"bytes"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//恢复文件块   （文件块的操作都是在内存中，待优化）
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
	switch params.StoragePolicy {
	case nodeClient.StoragePolicyCopy:
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
	case nodeClient.StoragePolicyEC:
		dataNodeNum := params.DataShardNum + params.ParityShardNum
		buffers := make([]*bytes.Buffer, dataNodeNum)
		for i := 0; i < dataNodeNum; i++ {
			buffers[i] = new(bytes.Buffer)
		}
		for index, shard := range params.Block.Shards {
			if shard.ShardID == params.ShardId {
				buffers[index] = nil
				continue
			}
			// tempPath := filepath.Join(config.TempDir, shard.Hash)
			// tempFile, err := os.Create(tempPath)
			// if err != nil {
			// 	logger.Logger.Error("创建EC临时数据文件失败", zap.Error(err))
			// 	ResponseErr(c, CodeServerBusy)
			// 	return
			// }
			err = config.DataNodeClient.ECDownloadShard(shard.Hash, shard.NodeURL, buffers[index])
			//err = config.DataNodeClient.ECDownloadShard(shard.Hash, shard.NodeURL, tempFile)
			if err != nil {
				logger.Logger.Warn("获取EC数据失败", zap.Error(err))
				buffers[index] = nil
			}
			// tempFile.Close()  //关闭句柄保存文件
			// tempFile, err = os.Open(tempPath)
			// if err != nil{
			// 	logger.Logger.Warn("打开EC临时数据文件失败", zap.Error(err))
			// 	buffers[index] = nil
			// }else{
			// 	buffers[index] =
			// }
		}
		shardPath := filepath.Join(config.ShardsDir, params.Block.Shards[params.ShardId].Hash)
		err = storagesdk.ReconstructMissShardFile(buffers, shardPath, params.DataShardNum, params.ParityShardNum, int(params.ShardId))
		if err != nil {
			logger.Logger.Error("恢复文件失败", zap.Error(err))
			ResponseErr(c, CodeFileError)
			return
		}
	}
	config.RecoveringShardHash.Delete(params.Block.Shards[params.ShardId].Hash)
	ResponseSuc(c, nil)
}
