package model

import "time"

// //文件表
// type Table_file struct {
// 	FileID    int       `json:"file_id"`
// 	FileSize  int64     `json:"file_size"`
// 	FileHash  string    `json:"file_hash"`
// 	FileURL   string    `json:"file_url"`
// 	Create_at time.Time `json:"create_at"`
// 	Update_at time.Time `json:"update_at"`
// }

// //用户文件表结构
// type UserFile struct {
// 	UserID    int64     `db:"user_id"`
// 	FileID    int       `db:"file_id"`
// 	FileName  string    `db:"file_name"`
// 	Create_at time.Time `db:"create_at"`
// 	Update_at time.Time `db:"update_at"`
// }

//文件meta信息
type FileMeta struct {
	FileID   int64  `db:"ID"`
	FileSize int64  `db:"file_size"`
	FileHash string `db:"file_hash"`
	FileURL  string `db:"file_url"`
}

//用户文件信息
type UserFileMeta struct {
	FileID   int64     `db:"file_id" json:"file_id"`
	FileSize int64     `db:"file_size" json:"file_size"`
	FileHash string    `db:"file_hash" json:"file_hash"`
	FileName string    `db:"file_name" json:"file_name"`
	FileURL  string    `db:"file_url" json:"-"`
	UpdateAt time.Time `db:"update_at" json:"update_at"`
}

//用户视频文件信息
type UserVideoMeta struct {
	FileID   int64     `db:"file_id" json:"file_id"`
	FileSize int64     `db:"file_size" json:"file_size"`
	FileHash string    `db:"file_hash" json:"file_hash"`
	FileName string    `db:"file_name" json:"file_name"`
	FaceUrl  string    `db:"face_url" json:"face_url"`
	UpdateAt time.Time `db:"update_at" json:"update_at"`
}

//用户更新文件信息
type UserUpdateFileInfo struct {
	FileID   int64  `json:"file_id" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

//用户删除文件信息
type UserDelFile struct {
	FileID int64  `json:"file_id,string" binding:"required"`
	Mod    string `json:"mod" binding:"required"`
}

//用户流式上传文件信息
type UserUploadFile struct {
	FileHash string `json:"file_hash" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

//用户分块上传文件初始化
type UserUploadMultiInit struct {
	FileHash string `json:"file_hash" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

//用户分块上传文件参数
type UserUploadMultiParts struct {
	ChunkIndex int    `form:"chunk_index" binding:"required"`
	FileName   string `form:"file_name" binding:"required"`
	UploadID   string `form:"upload_id" binding:"required"`
}

//用户分块上传进度查询
type UserUploadMultiProcess struct {
	UploadID string            `json:"upload_id" binding:"required"`
	Etags    map[string]string `json:"Etags"`
}

//用户完成分块上传参数
type UserUploadMultiComplete struct {
	FileName string `json:"file_name" binding:"required"`
	UploadID string `json:"upload_id" binding:"required"`
}

//用户分块上传文件缓存信息
type UserUploadMultiInfo struct {
	UploadID string
	FileName string
	FileHash string
	FileSize int64
}
