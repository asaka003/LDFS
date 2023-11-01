package sqlite

import (
	"LDFS/nameNode/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var LDB *gorm.DB

//初始化连接
func SqliteInit() (err error) {
	// 连接 SQLite 数据库
	LDB, err = gorm.Open(sqlite.Open(config.Sqlite.DbFile), &gorm.Config{})
	if err != nil {
		panic("无法连接到SqlLite数据库")
	}
	return
}
