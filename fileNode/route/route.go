package route

import (
	"LDFS/fileNode/controller"

	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	r := gin.Default()

	api := r.Group("/LDFS/")
	{
		api.POST("/replicasUploadShard", controller.ReplicasUploadShard)
		api.POST("/replicasDownloadShard", controller.ReplicasDownloadShard)

		api.POST("/ECuploadShard", controller.ECUploadShard)
		api.POST("/ECdownloadShard", controller.ECDownloadShard)

		api.POST("/recoverShard", controller.RecoverShard)
		api.POST("/getStorageInfo", controller.GetStorageInfo)

	}

	return r
}
