package model

import "time"

type ReqPutFileURL struct {
	Hash     string `json:"hash" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

type ReqUpdateFile struct {
	FileId   int64  `json:"file_id,string" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

type ReqDelFile struct {
	FileId int64 `json:"file_id,string" binding:"required"`
	//FileSystem string `json:"file_system" binding:"required"`
	Mod string `json:"mod"`
}

type ReqRecycleFile struct {
	FileId int64 `json:"file_id,string" binding:"required"`
}

type ReqCheckFileHash struct {
	FileHash string `json:"file_hash" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

type ReqUploadID struct {
	UploadID string `json:"upload_id" binding:"required"`
}

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

//用户分块上传文件初始化
type UploadMultiInit struct {
	FileKey string `json:"file_key" binding:"required"`
}

//用户分块上传文件参数
type UploadMultiParts struct {
	ChunkIndex int    `form:"chunk_index" binding:"required"`
	FileName   string `form:"file_key" binding:"required"`
	UploadID   string `form:"upload_id" binding:"required"`
}

//用户分块上传进度查询
type UploadMultiProcess struct {
	FileKey  string `json:"file_key" binding:"required"`
	UploadID string `json:"upload_id" binding:"required"`
	//Etags    map[string]string `json:"Etags"`
}

//用户完成分块上传参数
type UploadMultiComplete struct {
	FileKey  string `json:"file_key" binding:"required"`
	FileHash string `json:"file_hash" binding:"required"`
	UploadID string `json:"upload_id" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
}
