package mysql

import (
	"LDFS/nameNode/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" //使用mysql数据库都需要导入这个驱动才能识别运行
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Close() {
	_ = DB.Close()
}

func MysqlInit() error {
	log.Println("开始初始化mysql数据库连接...")
	var err error
	conn_Str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", //parseTime=true表示会将时间戳自动解析为具体时间存取数据库
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Host,
		config.Mysql.Port,
		config.Mysql.DB,
	)
	DB, err = sqlx.Connect("mysql", conn_Str)
	if err != nil {
		fmt.Println("链接mysql数据库错误!")
		return err
	}
	DB.SetMaxOpenConns(config.Mysql.MaxOpenConns)
	DB.SetMaxIdleConns(config.Mysql.MaxIdleConns)
	return nil
}
