package main

import (
	"LDFS/dataNode/config"
	"LDFS/dataNode/route"
	"LDFS/dataNode/util"
	"LDFS/logger"
	"LDFS/model"
	"flag"
	"log"
	"os"
	"os/signal"
)

const (
	DefaultHTTPAddr = "localhost:11000"
)

var httpAddr string
var joinNameNodeHaddr string
var nodeID string
var shardsDir string

func init() {
	err := config.ConfigInit()
	if err != nil {
		panic(err)
	}
	//初始化logger
	logger.InitLog()
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&joinNameNodeHaddr, "joinND", "", "Set join nameNode http address")
	flag.StringVar(&nodeID, "id", "", "Node ID.")
	flag.StringVar(&shardsDir, "shardsDir", "", "Set ths shards storage directory")
}

func main() {
	flag.Parse()

	if joinNameNodeHaddr == "" {
		log.Fatalf("need NameNodeAddr to Join")
	}

	if shardsDir != "" {
		config.ShardsDir = shardsDir
	}
	if err := os.MkdirAll(config.ShardsDir, 0700); err != nil {
		log.Fatalf("failed to create path for Shard storage: %s", err.Error())
	}
	// get dataNode info
	total, free, used, err := util.GetDiskUsageInfo()
	if err != nil {
		log.Fatalf("failed to get system disk usage info: %s", err.Error())
	}
	ShardsSize, err := util.GetDirectorySize(config.ShardsDir)
	if err != nil {
		log.Fatalf("failed to get shardsDir usage info: %s", err.Error())
	}
	dataNode := &model.DataNode{
		URL:                   "http://" + httpAddr,
		NodeName:              nodeID,
		NodeDiskSize:          total,
		NodeFileTotalSize:     ShardsSize,
		NodeDiskUsedSize:      used,
		NodeDiskAvailableSize: free,
	}

	if err := config.NameNodeClient.JoinDataNode(joinNameNodeHaddr, dataNode); err != nil {
		log.Fatalf("failed to join NameNode: %s", err.Error())
	}

	//获取当前所有NameNode列表信息
	config.NameNodeList, err = config.NameNodeClient.GetNameNodeListInfo("http://" + joinNameNodeHaddr)
	if err != nil {
		log.Fatalf("failed to get NameNode List: %s", err.Error())
	}

	//初始化路由信息
	r := route.SetRoute()
	go func() {
		r.Run(httpAddr)
	}()

	// We're up and running!
	log.Printf("DataNode %s started successfully, listening on http://%s", nodeID, httpAddr)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("DataNode exiting")
}
