package main

import (
	"LDFS/nameNode/config"
	"LDFS/nameNode/route"
)

func main() {
	if err := config.ConfigInit(); err != nil { //初始化配置信息
		panic(err)
	}
	//初始化路由信息
	r := route.SetRoute()
	r.Run(config.HttpApiServerHost)
}
