package model

import "time"

// type StorageFileList struct {
// 	FileList map[string]string
// }

//用户分块上传文件初始化
type UploadMultiInit struct {
	FileKey string `json:"file_key" binding:"required"`
}

//用户分块上传文件参数
type UploadMultiParts struct {
	ChunkIndex int    `form:"chunk_index" binding:"required"`
	FileName   string `form:"file_name" binding:"required"`
	UploadID   string `form:"upload_id" binding:"required"`
}

//用户分块上传进度查询
type UploadMultiProcess struct {
	UploadID string `json:"upload_id" binding:"required"`
	//Etags    map[string]string `json:"Etags"`
}

//用户完成分块上传参数
type UploadMultiComplete struct {
	//FileName string `json:"file_name" binding:"required"`
	UploadID string `json:"upload_id" binding:"required"`
}

//用户分块上传文件缓存信息
type UploadMultiInfo struct {
	UploadID string
	FileKey  string
	FileSize int64
}

type FileInfo struct {
	UUID       string    `json:"uuid"`
	Name       string    `json:"name"`
	FileKey    string    `json:"file_key"`
	DirType    string    `json:"dir_type"`
	CreateTime time.Time `json:"create_time"`
}

type DirInfo struct {
	DirPath   string      `json:"path"`
	FileMetas []*FileInfo `json:"files"`
}
