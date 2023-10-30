package sqlite

import (
	"LDFS/nameNode/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var LDB *gorm.DB

//初始化连接
func SqliteInit() (err error) {
	// 连接 SQLite 数据库
	LDB, err = gorm.Open(config.Sqlite.Dialect, config.Sqlite.DbFile)
	if err != nil {
		panic("无法连接到SqlLite数据库")
	}
	return
}

func Close() {
	LDB.Close()
}
