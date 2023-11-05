package controller

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/logger"
	"LDFS/dataNode/util"
	"LDFS/model"
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*副本冗余模式 Controller*/

//上传数据块
func ReplicasUploadShard(c *gin.Context) {
	shardHash := c.PostForm("hash")
	blockJson := c.PostForm("blockJson")
	block := &model.Block{}
	err := json.Unmarshal([]byte(blockJson), block)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	copyNo, err := strconv.ParseInt(c.PostForm("copyNo"), 10, 64)
	if err != nil || copyNo > int64(len(block.Shards)) {
		ResponseErr(c, CodeInvalidParam)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("读取请求体文件数据失败", zap.Error(err))
		return
	}
	ServerFileHash, err := util.MutiPartFileHash(fileHeader) //计算数据hash值
	if err != nil {
		logger.Logger.Error("计算文件hash值失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	if ServerFileHash != shardHash {
		//logger.Logger.Error("文件hash值不匹配", zap.Error(err))
		ResponseErr(c, CodeHashNotMatch)
		return
	}
	if strings.Contains(fileHeader.Filename, "..") { //过滤路径穿越
		ResponseErr(c, CodeFileNameInvalid)
		return
	}
	filePath := filepath.Join(config.ShardsDir, shardHash)
	c.SaveUploadedFile(fileHeader, filePath)

	//如不是最后一个副本，则继续副本冗余传递
	if copyNo < int64(len(block.Shards)-1) {
		buf := make([]byte, 0)
		buffer := bytes.NewBuffer(buf)
		src, err := fileHeader.Open()
		if err != nil {
			ResponseErr(c, CodeInvalidParam)
			return
		}
		defer src.Close()
		io.Copy(buffer, src)
		err = config.DataNodeClient.ReplicasUploadShard(shardHash, blockJson, buffer, block.Shards[copyNo+1].NodeURL, copyNo+1)
		if err != nil {
			ResponseErr(c, CodeServerBusy)
			return
		}
	}
	ResponseSuc(c, nil)
}

//下载数据块
func ReplicasDownloadShard(c *gin.Context) {
	shard := &model.Shard{}
	err := c.ShouldBindJSON(shard)
	if err != nil || strings.Contains(shard.Hash, "..") {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	shardPath := filepath.Join(config.ShardsDir, shard.Hash)
	c.File(shardPath)
}
