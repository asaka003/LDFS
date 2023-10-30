package route

import (
	"LDFS/fileNode/controller"

	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	r := gin.Default()

	r.GET("/getFileData/:type/:filename", controller.GetFileData) //获取文件块信息
	r.POST("/copyFileData", controller.CopyFileData)              //保存文件块信息

	// api2 := r.Group("/LDFS/nameNode/")
	// {
	// 	api2.GET("/download", controller.GetFile)
	// 	// api2.GET("/videoImg", controller.VideoImgHandler)
	// 	api2.POST("/upload", controller.SampleUploadFile)
	// 	api2.POST("/delFile", controller.DelFile)
	// }
	// api3 := r.Group("/LDFS/nameNode-multi/")
	// {
	// 	api3.POST("/MultiUploadInit", controller.UploadMultiInit)
	// 	api3.POST("/MultiUploadParts", controller.UploadMultiParts)
	// 	api3.POST("/MultiUploadComplete", controller.UploadMultiComplete)
	// 	api3.POST("/MultiUploadProcess", controller.UploadListParts)
	// }

	return r
}
