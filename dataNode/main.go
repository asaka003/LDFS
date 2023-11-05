package main

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/route"
)

func main() {
	config.ConfigInit()

	//启动http路由服务
	r := route.SetRoute()
	r.Run(config.AddrHttpLocalServer)
}
