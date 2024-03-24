package route

import (
	"LDFS/nameNode/controller"

	"github.com/gin-gonic/gin"
)

//文件上传反向代理(一致性hash算法单节点)
// func FileUploadReverseProxyHandler(c *gin.Context) {
// 	cluster := getClusterByPath(c.Request.URL.Path)
// 	if cluster == nil {
// 		http.Error(c.Writer, "404 not Found", http.StatusNotFound)
// 		return
// 	}
// 	backend, _ := cluster.Consistent.Get(c.GetHeader("FileHash"))
// 	backendURL, _ := url.Parse(backend)

// 	proxy := httputil.NewSingleHostReverseProxy(backendURL)

// 	proxy.ServeHTTP(c.Writer, c.Request)
// }

func SetRoute() *gin.Engine {
	r := gin.Default()

	api := r.Group("/LDFS/")
	{
		api.POST("/getAllFileKeys", controller.GetAllFileKeys)
		api.POST("/getFileMetaByFileKey/:fileKey", controller.GetFileMetaByFileKey)

		api.POST("/requestUploadFile", controller.RequestUploadFile)
		api.POST("/completeSampleUpload", controller.CompleteSampleUpload)

		api.POST("/getDataNodeListInfo", controller.GetDataNodeListInfo)
		// api.POST("/addDataNode", controller.AddDataNode)

		api.POST("/join", controller.JoinNameNodeHandler)
		api.POST("/joinDataNode", controller.JoinDataNodeHandler)
	}

	return r
}
