package main

import (
	"LDFS/nameNode/config"
	"LDFS/nameNode/raft"
	"LDFS/nameNode/route"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	DefaultHTTPAddr = "localhost:11000"
	DefaultRaftAddr = "localhost:12000"
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
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/LDFS/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
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

	config.RaftNode = raft.New()
	s := config.RaftNode
	s.RaftDir = raftDir
	s.MetaDir = metaDir
	s.RaftBind = raftAddr
	if err := s.Open(joinAddr == "", nodeID); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	// If join was specified, make the join request.
	if joinAddr != "" {
		if err := join(joinAddr, raftAddr, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
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
