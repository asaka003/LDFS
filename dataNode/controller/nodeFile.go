package controller

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/util"
	"LDFS/logger"
	"LDFS/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取存储占用信息
func GetStorageInfo(c *gin.Context) {
	dir := config.ShardsDir

	// 获取指定目录下的所有文件大小
	fileTotalSize, err := util.GetDirectorySize(dir)
	if err != nil {
		logger.Logger.Error("获取目录大小失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}
	// fmt.Println(dir)
	// fmt.Println(fileTotalSize)

	// // 获取程序运行的路径
	// executable, err := os.Executable()
	// if err != nil {
	// 	logger.Logger.Error("无法获取程序路径:", zap.Error(err))
	// 	return
	// }

	// // 获取磁盘的根目录
	// diskRoot := filepath.VolumeName(executable)

	// //获取磁盘使用大小
	// totalUsedSize, err := util.GetDirectorySize(diskRoot)
	// if err != nil {
	// 	logger.Logger.Error("获取目录大小失败", zap.Error(err))
	// 	ResponseErr(c, CodeServerBusy)
	// 	return
	// }

	// // 获取磁盘总大小
	// diskSize, err := util.GetSystemDiskSize()
	// if err != nil {
	// 	logger.Logger.Error("获取磁盘大小失败", zap.Error(err))
	// 	ResponseErr(c, CodeServerBusy)
	// 	return
	// }

	//获取磁盘存储信息
	total, free, used, err := util.GetDiskUsageInfo()
	if err != nil {
		logger.Logger.Error("获取磁盘存储信息失败", zap.Error(err))
		ResponseErr(c, CodeServerBusy)
		return
	}

	result := &model.DataNode{
		NodeDiskSize:          total,
		NodeFileTotalSize:     fileTotalSize,
		NodeDiskUsedSize:      used,
		NodeDiskAvailableSize: free,
	}

	// 构建响应
	ResponseSuc(c, result)
}
