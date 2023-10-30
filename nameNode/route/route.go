package route

import (
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

	// r.POST("/Login", controller.UserLogin) //已测试
	// r.POST("/Register", controller.UserRegister)
	// r.POST("/SendCode", controller.UserSignUpCode)
	// api := r.Group("/LDFS/nameNode", middleware.JWTAuthMiddleware, middleware.RateLimitMiddleware(time.Microsecond, 1000))
	// {
	// 	api.GET("/UserGetAllFilesList", controller.GetUserAllFilesList) //已测试
	// 	api.GET("/videoImg", controller.VideoImgHandler)                //已测试
	// 	api.GET("/download", controller.DownloadHandler)                //已测试
	// 	api.GET("/videoInfo", controller.GetVideoInfo)                  //已测试
	// 	api.GET("/videoPlay", controller.VideoPlayer)                   //已测试

	// 	api.POST("/UserUpdateFileInfo", controller.UpdateUserFileInfo)
	// 	api.POST("/delFile", controller.DelUserFile) //已测试
	// 	api.POST("/UserRecycleFile", controller.RecycleFile)
	// 	api.POST("/QuickUpload", controller.QuickUpload)

	// 	api.POST("/upload", controller.UserSimpleUploadFile)                 //已测试普通文件(未测试视频文件)
	// 	api.POST("/multiUploadInit", controller.UesrUploadMultiInit)         //已测试
	// 	api.POST("/multiUploadPart", controller.UserUploadMultiParts)        //已测试
	// 	api.POST("/multiUploadComplete", controller.UserUploadMultiComplete) //已测试
	// 	api.POST("/listParts", controller.UserUploadListParts)               //已测试
	// }

	return r
}
