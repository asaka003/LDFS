package db

import (
	"LDFS/model"
	"LDFS/nameNode/config"
)

var DB AbstractDB

type AbstractDB interface {
	SaveFileMetaInfo(fileKey string, fileMeta *model.FileMetadata) (err error)
	GetAllFileKeys() (err error)
	GetFileMetaByFileKey(fileKey string) (fileMeta *model.FileMetadata, err error)
	DelFileMetaByFileKey(fileKey string) (err error)
	UpdateFileMetaByFileKey(fileKey string, fileMeta *model.FileMetadata) (err error)
}

//初始化DB服务
func InitDB() {
	switch config.SystemDB {
	case "sqlite":
	case "mysql":
	case "redis":
	}
}
