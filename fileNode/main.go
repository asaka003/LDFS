package main

import (
	"LDFS/fileNode/config"
	"LDFS/fileNode/route"
	"LDFS/fileNode/web_pkg/snowflake"
)

func init() {
	//初始化配置文件
	if err := config.ViperInit(); err != nil {
		panic(err)
	}

}

func main() {

	//初始化文件存储路径
	config.FileDirInit()

	//初始化雪花算法
	snowflake.Init("2022-01-01", 1)

	//启动http路由服务
	r := route.SetRoute()
	r.Run(config.AddrHttpLocalServer)
}
