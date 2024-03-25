package main

import (
	"LDFS/logger"
	"LDFS/nameNode/config"
	"LDFS/nameNode/raft"
	"LDFS/nameNode/route"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

const (
	DefaultHTTPAddr = "http://" + "localhost:11000"
	DefaultRaftAddr = "http://" + "localhost:12000"
)

var httpAddr string
var raftAddr string
var joinAddr string
var nodeID string
var metaDir string

func init() {
	if err := config.ConfigInit(); err != nil { //初始化配置信息
		panic(err)
	}
	//初始化logger
	logger.InitLog()
	//flag.BoolVar(&inmem, "inmem", false, "Use in-memory storage for Raft")
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID. If not set, same as Raft bind address")
	flag.StringVar(&metaDir, "metadir", config.FileMetaDir, "file meta directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Int64Var(&config.ECDataShardNum, "d", config.ECDataShardNum, "set the number of RS data shard")
	flag.Int64Var(&config.ECParityShardNum, "p", config.ECDataShardNum, "set the number of RS parity shard")
	flag.Int64Var(&config.CopyReplicasNum, "r", config.CopyReplicasNum, "set the number of replicas")
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	if nodeID == "" {
		nodeID = raftAddr
	}

	// Ensure Raft storage exists.
	raftDir := flag.Arg(0)
	os.RemoveAll(raftDir)
	if raftDir == "" {
		log.Fatalln("No Raft storage directory specified")
	}
	if err := os.MkdirAll(raftDir, 0700); err != nil {
		log.Fatalf("failed to create path for Raft storage: %s", err.Error())
	}

	if metaDir == "" {
		log.Fatalln("No meta storage directory specified")
	}
	if err := os.MkdirAll(metaDir, 0700); err != nil {
		log.Fatalf("failed to create path for Meta storage: %s", err.Error())
	}

	err := raft.New(raftDir, metaDir, raftAddr, joinAddr, nodeID)
	if err != nil {
		panic(err)
	}

	//初始化路由信息
	r := route.SetRoute()
	go func() {
		r.Run(httpAddr)
	}()

	// We're up and running!
	log.Printf("NameNode %s started successfully, listening on http://%s", nodeID, httpAddr)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("NameNode exiting")
}
