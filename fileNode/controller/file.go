package controller

import (
	"LDFS/fileNode/config"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

//文件数据冗余复制(未进行认证，不安全)
func CopyFileData(c *gin.Context) {
	var dst string
	SaveType := c.PostForm("type")
	if SaveType == "" {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	header, err := c.FormFile("file")
	if err != nil {
		ResponseErr(c, CodeInvalidParam)
		return
	}
	switch SaveType {
	case "shard":
		dst = filepath.Join(config.FileShardDir, header.Filename)
	case "file":
		dst = filepath.Join(config.FilesDir, header.Filename)
	default:
		ResponseErr(c, CodeInvalidParam)
		return
	}
	err = c.SaveUploadedFile(header, dst)
	if err != nil {
		ResponseErr(c, CodeServerBusy)
		return
	}
	ResponseSuc(c, nil)
}

//文件数据冗余获取(未进行认证，不安全)
func GetFileData(c *gin.Context) {
	DataType := c.Param("type")
	filename := c.Param("filename")
	var dir string
	switch DataType {
	case "shard":
		dir = config.FileShardDir
	case "file":
		dir = config.FilesDir
	default:
		ResponseErr(c, CodeInvalidParam)
		return
	}
	c.File(filepath.Join(dir, filename))
}

// //获取文件完整数据(未进行认证,不安全)
// func GetFile(c *gin.Context) {
// 	FileKey := c.Query("FileKey")
// 	UUID, err := logic.GetUUIDByFileKey(FileKey)
// 	if err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	filePath, err := util.ReconstructFile(UUID)
// 	if err != nil {
// 		log.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(FileKey)))
// 	c.File(filePath) //自动处理Range
// }

// //删除文件数据
// func DelFile(c *gin.Context) {
// 	FileKey := c.Param("FileKey")
// 	err := logic.DelFile(FileKey) //删除对应的meta信息和数据块信息(--------这里没有实际删除磁盘中的数据库信息----------)
// 	if err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, nil)
// }

// //获取存储文件列表
// func GetStorageFileList(c *gin.Context) {
// 	dirPath := c.Param("dir")
// 	dirInfo, err := logic.GetDirInfo(dirPath)
// 	if err != nil {
// 		log.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, dirInfo)
// }

// //简单文件上传
// func SampleUploadFile(c *gin.Context) {
// 	fileKey := c.PostForm("FileKey")
// 	if fileKey == "" {
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//解析文件信息
// 	header, err := c.FormFile("file")
// 	if err != nil {
// 		ResponseErr(c, CodeInvalidParam)
// 		log.Println(err)
// 		return
// 	}
// 	UUID := uuid.New().String()
// 	filePath := filepath.Join(config.FilesDir, UUID)
// 	//流式处理文件，保存文件到硬盘中
// 	err = c.SaveUploadedFile(header, filePath)
// 	if err != nil {
// 		log.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, nil)

// 	err = util.DistributeFileToNodes(UUID, filePath, fileKey)
// 	if err != nil {
// 		log.Println("纠删码冗余文件数据失败", err)
// 	}
// }

//初始化分块上传文件信息
// func UploadMultiInit(c *gin.Context) {
// 	params := new(model.UploadMultiInit)
// 	if err := c.ShouldBind(params); err != nil {
// 		log.Println("解析参数失败", err)
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	uploadID := util.GenerateUploadID(params.FileKey) //生成uploadID值
// 	//创建分块上传文件缓存信息
// 	if err := logic.CreateMultiInfo(uploadID, params.FileKey); err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, uploadID)
// }

//分块上传文件
// func UploadMultiParts(c *gin.Context) {
// 	params := new(model.UploadMultiParts)
// 	chunkIndex, err := strconv.ParseInt(c.PostForm("ChunkIndex"), 10, 64)
// 	if err != nil {
// 		log.Println(err)
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	params.ChunkIndex = int(chunkIndex)
// 	params.UploadID = c.PostForm("UploadID")
// 	params.FileName = c.PostForm("FileKey")

// 	if params.UploadID == "" || params.FileName == "" || !util.IsUploadID(params.UploadID) {
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//读取上传分块文件信息
// 	file_head, err := c.FormFile("file")
// 	if err != nil {
// 		log.Println("获取文件失败", zap.Error(err))
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//查询uploadID是否存在
// 	ok, err := logic.IsUploadIDExist(params.UploadID)
// 	if err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	if !ok {
// 		ResponseErr(c, CodeUploadIDNotFound)
// 		return
// 	}

// 	//保存分块文件信息
// 	file_dir := filepath.Join(config.FilesDir, params.UploadID)
// 	util.CreateDir(file_dir)
// 	filePath := filepath.Join(config.FilesDir, params.UploadID, strconv.FormatInt(chunkIndex, 10))
// 	err = c.SaveUploadedFile(file_head, filePath)
// 	if err != nil {
// 		fmt.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ETag, err := util.MutiFileHash(file_head)
// 	if err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	//更新分块文件缓存信息
// 	if err := logic.UpdateMultiInfo(params.UploadID, params.ChunkIndex, ETag, file_head.Size); err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, nil)
// }

//分块上传进度查询
// func UploadListParts(c *gin.Context) {
// 	var err error
// 	params := new(model.UploadMultiProcess)
// 	if err = c.ShouldBindJSON(params); err != nil {
// 		log.Println("解析参数失败", err)
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//读取分块文件传ETag
// 	Etags, err := logic.GetAllETags(params.UploadID)
// 	if err != nil {
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	parts := []string{}
// 	for index := range Etags.Etags {
// 		parts = append(parts, index)
// 	}
// 	ResponseSuc(c, parts)
// }

// //分块文件上传完成提交
// func UploadMultiComplete(c *gin.Context) {
// 	params := new(model.UploadMultiComplete)
// 	if err := c.ShouldBindJSON(params); err != nil {
// 		log.Println("解析参数失败", zap.Error(err))
// 		ResponseErr(c, CodeInvalidParam)
// 		return
// 	}
// 	//根据uploadID获取分块上传文件信息
// 	multiInfo, err := logic.GetUploadMultiInfo(params.UploadID)
// 	if err != nil {
// 		fmt.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	UUID := uuid.New().String()

// 	tmpMultiDir := filepath.Join(config.FilesDir, params.UploadID)
// 	newfilePath := filepath.Join(config.FilesDir, UUID)
// 	if err := util.MergeMultiFile(tmpMultiDir, newfilePath); err != nil {
// 		fmt.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	// //提取文件封面信息
// 	// mimetype := mime.TypeByExtension(filepath.Ext(params.FileName))
// 	// // 判断 MIME 类型是否是视频文件
// 	// if strings.HasPrefix(mimetype, "video/") {
// 	// 	img_path := filepath.Join(config.VideoFaceImgDir, UUID+".jpg")
// 	// 	InputPath := newfilePath

// 	// 	_, err = VideoServerCli.GenerateVideoImg(context.TODO(), &proto.ReqVideoInfo{ //提取视频封面信息
// 	// 		Option:         kafka.OptionImg, //可以删除省略
// 	// 		StorageType:    kafka.StorageTypeLocal,
// 	// 		FileKey:        "-",
// 	// 		VideoImgPath:   img_path,
// 	// 		VideoInputPath: InputPath,
// 	// 	})
// 	// 	if err != nil {
// 	// 		log.Println("提交分封面截取任务失败", err)
// 	// 	}
// 	// 	_, err = VideoServerCli.GenerateLowVideoResolution(context.TODO(), &proto.ReqVideoInfo{
// 	// 		Option:         kafka.OptionLowResolution, //可以删除省略
// 	// 		StorageType:    kafka.OptionLowResolution,
// 	// 		FileKey:        "-",
// 	// 		VideoInputPath: InputPath,
// 	// 	})
// 	// 	if err != nil {
// 	// 		log.Println("提交低分辨率视频处理任务失败", err)
// 	// 	}
// 	// }

// 	//删除用户分块上传文件缓存信息
// 	if err = logic.ClearUploadMultiData(params.UploadID); err != nil {
// 		fmt.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}

// 	//将合并完成的文件冗余到其他节点中(纠删码模式)
// 	err = util.DistributeFileToNodes(UUID, newfilePath, multiInfo.FileKey)
// 	if err != nil {
// 		log.Println(err)
// 		ResponseErr(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuc(c, nil)
// }

// //处理视频文件
// func VideoHandler(c *gin.Context) {
// 	//fileKey := c.Param("FileKey")
// }
