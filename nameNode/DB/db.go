package db

import (
	"LDFS/model"
	"LDFS/nameNode/config"
)

var DB AbstractDB

type AbstractDB interface {
	SaveSampleUploadInfo(fileKey string, list []*model.Shard) (err error)
}

//初始化DB服务
func InitDB() {
	switch config.SystemDB {
	case "sqlite":
	case "mysql":
	case "redis":
	}
}
