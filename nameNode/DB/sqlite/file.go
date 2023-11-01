package sqlite

import (
	"LDFS/model"

	"gorm.io/gorm"
)

type SqlLiteDB struct {
	db *gorm.DB
}

//读取指定目录所有文件列表
func (sqlite *SqlLiteDB) GetAllFileKeys() (err error) {
	return
}

//读取指定fileKey的文件Meta信息
func (sqlite *SqlLiteDB) GetFileMetaByFileKey(fileKey string) (fileMeta *model.FileMetadata, err error)

//添加文件Meta信息
func (sqlite *SqlLiteDB) SaveFileMetaInfo(fileKey string, fileMeta *model.FileMetadata) (err error)

//更新文件Meta信息
func (sqlite *SqlLiteDB) UpdateFileMetaByFileKey(fileKey string, fileMeta *model.FileMetadata) (err error)

//删除文件Meta信息
func (sqlite *SqlLiteDB) DelFileMetaByFileKey(fileKey string) (err error)
