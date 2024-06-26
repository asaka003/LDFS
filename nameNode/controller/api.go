package controller

import (
	"LDFS/logger"
	"LDFS/model"
	"LDFS/nameNode/config"
	"LDFS/nameNode/raft"
	"LDFS/nameNode/util"
	"LDFS/nodeClient"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取所有文件信息列表
func GetAllFileKeys(c *gin.Context) {
	dir, err := os.Open(config.FileMetaDir)
	if err != nil {
		ResponseErr(c, CodeNotFoundFile)
		return
	}
	defer dir.Close()

	fileList, err := dir.Readdir(-1)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}

	fileInfoList := make([]*model.FileInfo, 0)
	for _, f := range fileList {
		if f.IsDir() {
			continue
		}
		//读取文件meta信息
		path := filepath.Join(config.FileMetaDir, f.Name())
		fileMeta, err := util.GetFileMetaInFile(path)
		if err != nil {
			continue
		}
		fileInfoList = append(fileInfoList, &model.FileInfo{
			FileKey: fileMeta.FileKey,
			Size:    fileMeta.FileSize,
		})
	}
	ResponseSuc(c, fileInfoList)
}

//根据fileKey获取文件Meta信息
func GetFileMetaByFileKey(c *gin.Context) {
	fileKey := c.Param("fileKey")
	if strings.Contains(fileKey, "..") {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	fileMeta, err := util.GetFileMetaInFile(fileKey)
	if err != nil {
		//fmt.Println(fileKey)
		logger.Logger.Error("读取fileMeta文件信息失败", zap.Error(err))
		ResponseErr(c, CodeNotFoundFile)
		return
	}
	ResponseSuc(c, fileMeta)
}

/*
	处理不同策略模式的文件上传，仅初始化文件的meta信息
	- 如果是副本冗余模式，则会根据nameNode设置的block_size大小进行分配dataNode
	- 如果是纠删码模式，则会根据nameNode设置的block_size划分不同的文件块
	随后对不同的文件块进行数据块拆分，给出数据块和校验块的dataNode存储地址
*/
func RequestUploadFile(c *gin.Context) {
	params := new(model.RequestUploadFileParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	fileMeta := &model.FileMetadata{}
	fileMeta.FileKey = params.FileKey
	fileMeta.StoragePolicy = params.StoragePolicy
	fileMeta.FileSize = params.FileSize
	fileMeta.CreateTime = time.Now()

	result := &model.RequestUploadFileResponse{
		FileMeta: fileMeta,
	}

	blockNum := int64(math.Ceil(float64(params.FileSize) / float64(params.BlockSize)))
	//获取所有能够存储block的DataNode列表
	//dataNodeLen := len(config.DataNodeList)
	availableDataNodeList := make([]model.DataNode, 0)
	for _, dataNode := range raft.RaftNodeClient.GetDataNodeList() {
		if int64(dataNode.NodeDiskAvailableSize)-config.RemainSize > params.FileSize {
			availableDataNodeList = append(availableDataNodeList, *dataNode)
		}
	}
	if len(availableDataNodeList) == 0 { //所有的dataNode都存储满，返回错误信息
		logger.Logger.Error("Disk is full")
		ResponseErr(c, CodeDiskIsFull)
		return
	}

	var selectDataNodeNum int
	switch params.StoragePolicy {
	case nodeClient.StoragePolicyCopy:
		selectDataNodeNum = int(config.CopyReplicasNum)
	case nodeClient.StoragePolicyEC:
		selectDataNodeNum = int(config.ECDataShardNum + config.ECParityShardNum)
		fileMeta.DataShards = int(config.ECDataShardNum)
		fileMeta.ParityShards = int(config.ECParityShardNum)
	default:
		ResponseErr(c, CodeInvalidParam)
		return
	}

	//根据DataNode剩余容量选择不重复的DataNode，如果选取的数量大于实际运行的dataNode数量，则会重复选取
	for blockId := 0; blockId < int(blockNum); blockId++ {
		//对所有可用的DataNode列表中的剩余空间进行排序
		sort.Slice(availableDataNodeList, func(i, j int) bool {
			return availableDataNodeList[i].NodeDiskAvailableSize > availableDataNodeList[j].NodeDiskAvailableSize
		})

		selectDataNodeList := make([]model.DataNode, 0)

		for i := 0; i < selectDataNodeNum; i++ {
			if i >= len(availableDataNodeList) { //如果需要的数量大于可以DataNode，则单个DataNode存储多个副本
				availableDataNodeList[i-len(availableDataNodeList)].NodeDiskUsedSize += params.BlockSize
				availableDataNodeList[i-len(availableDataNodeList)].NodeDiskAvailableSize -= params.BlockSize
				selectDataNodeList = append(selectDataNodeList, availableDataNodeList[i-len(availableDataNodeList)])
			} else {
				availableDataNodeList[i].NodeDiskUsedSize += params.BlockSize
				availableDataNodeList[i].NodeDiskAvailableSize -= params.BlockSize
				selectDataNodeList = append(selectDataNodeList, availableDataNodeList[i])
			}
		}
		shards := make([]*model.Shard, 0)
		for i, dataNode := range selectDataNodeList {
			shards = append(shards, &model.Shard{
				ShardID: int64(i),
				NodeURL: dataNode.URL,
			})
		}

		//计算每个block对应的实际大小
		blockSize := params.FileSize - int64(blockId)*params.BlockSize
		if blockSize > params.BlockSize {
			blockSize = params.BlockSize
		}
		fileMeta.Blocks = append(fileMeta.Blocks, &model.Block{
			BlockId:   blockId,
			BlockSize: blockSize,
			Shards:    shards,
		})
	}

	//保存meta信息到文件中
	err = util.SaveFileMetaInFile(fileMeta)
	if err != nil {
		logger.Logger.Error("failed to save meta", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuc(c, result)
}

//完成简单文件上传请求
func CompleteSampleUpload(c *gin.Context) {
	params := new(model.CompleteSampleUploadParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	fileMeta, err := util.GetFileMetaInFile(params.FileKey)
	if err != nil {
		logger.Logger.Error("读取fileMeta文件信息失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	//更新hash值
	if fileMeta.StoragePolicy == nodeClient.StoragePolicyEC {
		for i, block := range params.FileMeta.Blocks {
			fileMeta.Blocks[i].Hash = block.Hash
			for j, shard := range block.Shards {
				fileMeta.Blocks[i].Shards[j].Hash = shard.Hash
			}
		}
	}

	fileMeta.Status = "success"
	err = util.SaveFileMetaInFile(fileMeta)
	if err != nil {
		logger.Logger.Error("保存fileMeta文件信息失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuc(c, nil)
}

//删除文件信息
func DeleteFile(c *gin.Context) {
	fileKey := c.Param("fileKey")
	if strings.Contains(fileKey, "..") {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	err := util.DeleteFileMeta(fileKey)
	if err != nil {
		logger.Logger.Error("删除fileMeta文件信息失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuc(c, nil)
}

//修改文件名
func UpdateFileName(c *gin.Context) {}

//更新文件meta信息
func UpdateFileMeta(c *gin.Context) {
	fileMeta := new(model.FileMetadata)
	err := c.ShouldBindJSON(fileMeta)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	//fmt.Println(fileMeta.Blocks[0].Hash)
	err = util.UpdateFileMeta(fileMeta)
	if err != nil {
		logger.Logger.Error("更新文件meta信息失败", zap.Error(err))
		if err.Error() == raft.ErrNotLeader {
			ResponseErr(c, CodeNotLeader)
		} else {
			ResponseErr(c, CodeServerBusy)
		}
		return
	}
	ResponseSuc(c, nil)
}
