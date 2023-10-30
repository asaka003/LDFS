package controller

import (
	"LDFS/fileNode/logger"
	"LDFS/model"
	"LDFS/nameNode/config"
	"LDFS/nameNode/logic"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取所有文件元信息列表

//根据指定的ID获取文件元信息

//初始化分块上上传
func InitMultiUpload(c *gin.Context) {
	params := new(model.InitUploadParam)
	err := c.ShouldBindJSON(params)
	if err != nil || strings.Contains(params.FileHash, "..") { //过滤路径穿越
		ResponseErr(c, CodeInvalidParam)
		return
	}

	//创建对应的追加写入文件
	path := filepath.Join(config.MultiUploadDir, params.FileHash)
	file, err := os.Create(path)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	defer file.Close()
	ResponseSuc(c, nil)
}

//上传文件分块
func UploadMultiPart(c *gin.Context) {

}

//完成分块上传
func CompeleteMultiUpload(c *gin.Context) {

}

//查询文件上传进度
func CheckMultiProgress(c *gin.Context) {

}

//请求文件下一个分块上传DataNodes地址

//简单上传文件
func GetSampleUploadNodeList(c *gin.Context) {
	res := model.SampleUploadList{}
	ResponseSuc(c, res)
}

//保存简单上传文件DataNode存储列表信息
func SendSampleUploadInfo(c *gin.Context) {
	params := new(model.SampleUploadInfo)
	err := c.ShouldBindJSON(params)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	err = logic.SaveSampleUploadInfo(params.FileKey, params.Shards)
	if err != nil {
		logger.Logger.Error("保存简单上传文件失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuc(c, nil)
}

//请求上传文件
func RequestUploadFile(c *gin.Context) {
	params := new(model.RequestUploadFileParams)
	err := c.ShouldBindJSON(params)
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	result := &model.RequestUploadFileResponse{}
	switch params.StoragePolicy {
	case config.StoragePolicyCopy:
		blockNum := int64(math.Ceil(float64(params.FileSize) / float64(params.BlockSize)))
		//获取所有能够存储block的DataNode列表
		dataNodeLen := len(config.DataNodeList)
		for i := 0; i < int(blockNum); i++ {
			//随机选择三个不一样的DataNode

		}

	case config.StoragePolicyEC:

	}

	ResponseSuc(c, result)
}
