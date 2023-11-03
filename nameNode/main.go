package main

import (
	"LDFS/nameNode/config"
	"LDFS/nameNode/route"
	"fmt"

	"github.com/spf13/viper"
)

var (
	AddrHttpApiServer string
)

func init() {
	if err := config.ConfigInit(); err != nil { //初始化配置信息
		panic(err)
	}
	fmt.Println("初始化配置信息成功")
	AddrHttpApiServer = viper.GetString("Servers.ApiServer.IP")
}

func main() {

	//初始化路由信息
	r := route.SetRoute()
	r.Run(AddrHttpApiServer)
}
