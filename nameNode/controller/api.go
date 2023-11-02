package controller

import (
	"LDFS/fileNode/logger"
	"LDFS/model"
	"LDFS/nameNode/config"
	"LDFS/nameNode/util"
	"bytes"
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"

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
	metaPath := filepath.Join(config.FileMetaDir, util.BytesHash([]byte(fileKey))+".json")
	fileMeta, err := util.GetFileMetaInFile(metaPath)
	if err != nil {
		logger.Logger.Error("读取fileMeta文件信息失败", zap.Error(err))
		return
	}
	ResponseSuc(c, fileMeta)
}

//请求上传文件
func RequestUploadFile(c *gin.Context) {
	params := new(model.RequestUploadFileParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	fileMeta := &model.FileMetadata{}
	result := &model.RequestUploadFileResponse{
		FileMeta: fileMeta,
	}

	blockNum := int64(math.Ceil(float64(params.FileSize) / float64(params.BlockSize)))
	//获取所有能够存储block的DataNode列表
	//dataNodeLen := len(config.DataNodeList)
	availableDataNodeList := make([]model.DataNode, 0)
	for _, dataNode := range config.DataNodeList {
		if dataNode.NodeDiskSize-dataNode.NodeDiskUsedSize-config.RemainSize > params.FileSize {
			availableDataNodeList = append(availableDataNodeList, dataNode)
		}
	}
	if len(availableDataNodeList) == 0 { //所有的dataNode都存储满，返回错误信息
		ResponseErr(c, CodeDiskIsFull)
		return
	}

	var selectDataNodeNum int

	switch params.StoragePolicy {
	case config.StoragePolicyCopy:
		selectDataNodeNum = int(config.CopyReplicasNum)
	case config.StoragePolicyEC:
		selectDataNodeNum = int(config.ECDataShardNum + config.ECParityShardNum)
	default:
		ResponseErr(c, CodeInvalidParam)
		return
	}

	//根据DataNode剩余容量选择不重复的DataNode，如果选取的数量大于实际运行的dataNode数量，则会重复选取
	for bolckId := 0; bolckId < int(blockNum); bolckId++ {
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
		URLs := make([]string, 0)
		for _, dataNode := range selectDataNodeList {
			URLs = append(URLs, dataNode.URL)
		}
		fileMeta.Shards = append(fileMeta.Shards, &model.Shard{
			ShardID:  bolckId,
			NodeURLs: URLs,
		})
	}

	//保存meta信息到文件中
	metaJson, err := json.Marshal(fileMeta)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	path := filepath.Join(config.FileMetaDir, util.BytesHash([]byte(params.FileKey))+".json")
	_, err = os.Stat(path)
	if err == nil {
		ResponseErr(c, CodeFileExist)
		return
	}

	//创建文件目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	file, err := os.Create(path)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	defer file.Close()
	io.Copy(file, bytes.NewBuffer(metaJson))
	ResponseSuc(c, result)
}
