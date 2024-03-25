package controller

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/util"
	"LDFS/logger"
	"LDFS/model"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*纠删码模式*/
//上传文件数据块
func ECUploadShard(c *gin.Context) {
	fileHash := c.PostForm("hash")
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
	if ServerFileHash != fileHash {
		//logger.Logger.Error("文件hash值不匹配", zap.Error(err))
		ResponseErr(c, CodeHashNotMatch)
		return
	}
	if strings.Contains(fileHeader.Filename, "..") { //过滤路径穿越
		ResponseErr(c, CodeFileNameInvalid)
		return
	}
	filePath := filepath.Join(config.ShardsDir, fileHash)
	c.SaveUploadedFile(fileHeader, filePath)
}

//下载文件数据块
func ECDownloadShard(c *gin.Context) {
	params := new(model.DownloadShardParam)
	if err := c.ShouldBindJSON(params); err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	if strings.Contains(params.Hash, "..") { //过滤路径穿越
		ResponseErr(c, CodeFileNameInvalid)
		return
	}
	filePath := filepath.Join(config.ShardsDir, params.Hash)
	c.File(filePath)
}
