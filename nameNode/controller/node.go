package controller

import (
	"LDFS/model"
	"LDFS/nameNode/raft"

	"github.com/gin-gonic/gin"
)

/* 本模块负责节点之间的通信地址查询 */

//查询所有DataNode的信息
func GetDataNodeListInfo(c *gin.Context) {
	ResponseSuc(c, raft.RaftNodeClient.GetDataNodeList())
}

//加入NameNode节点
func JoinNameNodeHandler(c *gin.Context) {
	params := new(model.ParamJoin)
	if err := c.ShouldBindJSON(params); err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	raft.RaftNodeClient.Join(params.ID, params.RaftAddr)
	ResponseSuc(c, nil)
}

//加入DataNode节点
func JoinDataNodeHandler(c *gin.Context) {
	params := new(model.ParamJoinDataNode)
	if err := c.ShouldBindJSON(params); err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	raft.RaftNodeClient.AddDataNode(params.DataNodeInfo)
	ResponseSuc(c, nil)
}

//动态加入DataNode请求(没有密码认证，不安全)
// func AddDataNode(c *gin.Context) {
// 	param := new(model.DataNode)
// 	if err := c.ShouldBindJSON(param); err != nil {
// 		logger.Logger.Error("解析动态节点加入参数失败", zap.Error(err))
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//查询DataNode 存储磁盘信息
// 	config.DataNodeUrls = append(config.DataNodeUrls, param.URL)
// 	//请求DataNode 磁盘存储情况
// 	dataNode, err := config.DataNodeClient.GetStorageInfo(param.URL)
// 	if err != nil {
// 		logger.Logger.Error("获取DataNode节点信息失败", zap.Error(err))
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	config.DataNodeList = append(config.DataNodeList, dataNode)
// 	// 更新配置数据
// 	viper.Set("DataNodes.List", config.DataNodeList)
// 	// 保存到文件
// 	if err := viper.WriteConfigAs("config/ServerAddr.yaml"); err != nil {
// 		logger.Logger.Error("写入配置文件失败", zap.Error(err))
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, nil)
// }

// //根据名称查询节点服务地址端口
// func GetAddrByServerName(c *gin.Context) {
// 	param := new(model.FileNode)
// 	if err := c.ShouldBindJSON(param); err != nil {
// 		logger.Logger.Error("文件节点服务地址查询参数解析失败", zap.Error(err))
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	addr, ok := DataNodeUrl[param.ServerName]
// 	if !ok {
// 		ResponseErr(c, CodeNodeNotFound)
// 		return
// 	}
// 	strs := strings.SplitN(addr, ":", 2)
// 	param.IP = strs[0]
// 	param.Port = strs[1]
// 	ResponseSuc(c, param)
// }
