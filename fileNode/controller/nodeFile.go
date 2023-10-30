package controller

import (
	"LDFS/fileNode/config"
	"LDFS/fileNode/logger"
	"LDFS/fileNode/util"
	"LDFS/model"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取存储占用信息
func GetStorageInfo(c *gin.Context) {
	dir := config.FileShardDir

	// 获取指定目录下的所有文件大小
	fileTotalSize, err := util.GetDirectorySize(dir)
	if err != nil {
		logger.Logger.Error("获取目录大小失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 获取程序运行的路径
	executable, err := os.Executable()
	if err != nil {
		logger.Logger.Error("无法获取程序路径:", zap.Error(err))
		return
	}
	// 获取磁盘的根目录
	diskRoot := filepath.VolumeName(executable)

	//获取磁盘使用大小
	totalUsedSize, err := util.GetDirectorySize(diskRoot)
	if err != nil {
		logger.Logger.Error("获取目录大小失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	// 获取磁盘总大小
	diskSize, err := util.GetSystemDiskSize()
	if err != nil {
		logger.Logger.Error("获取磁盘大小失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	result := model.DataNode{
		NodeDiskSize:      diskSize,
		NodeFileTotalSize: fileTotalSize,
		NodeDiskUsedSize:  totalUsedSize,
	}

	// 构建响应
	ResponseSuc(c, result)
}

//恢复文件块数据
func RecoverShard(c *gin.Context) {

}

//复制其他节点的文件数据块(冗余副本策略)
func CopyFileShardFromOtherNode(c *gin.Context) {

}
