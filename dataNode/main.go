package main

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/route"
)

func main() {
	err := config.ConfigInit()
	if err != nil {
		panic(err)
	}

	//启动http路由服务
	r := route.SetRoute()
	r.Run(config.AddrHttpLocalServer)
}
