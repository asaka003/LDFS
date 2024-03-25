package raft

import (
	"LDFS/logger"
	"LDFS/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"go.uber.org/zap"
)

const ErrNotLeader = "err not leader"

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second

	CommandCreateFileMeta = "create-meta"
	CommandDeleteFileMeta = "delete-meta"
	CommandUpdateFileMeta = "update-meta"

	CommandAddDataNode = "add-dataNode"

	CommandAddNameNodeHaddr = "add-name-node-haddr"
)

type FileMeta model.FileMetadata

type RaftNode struct {
	RaftDir     string     //raft节点状态目录
	MetaDir     string     //文件meta信息目录
	RaftBind    string     //raft协议交互地址
	mu          sync.Mutex //锁用于控制文件meta操作的并发
	fileMutexs  map[string]*sync.RWMutex
	raft        *raft.Raft //raft算法核心组件
	logger      *log.Logger
	DataNodeSet map[string]*model.DataNode //管理的dataNode列表
	NameNodeSet map[string]*model.NameNode
	//DataNodeList []*model.DataNode
}

//获取文件meta信息
func (node *RaftNode) GetFileMeta(key string) (meta *FileMeta, err error) {
	if strings.Contains(key, "..") {
		return nil, errors.New("invalid param")
	}
	node.mu.Lock()
	fm, ok := node.fileMutexs[key]
	if !ok {
		node.fileMutexs[key] = &sync.RWMutex{}
		fm = node.fileMutexs[key]
	}
	node.mu.Unlock()
	fm.RLock()
	defer fm.RUnlock() //位置可优化

	path := filepath.Join(node.MetaDir, key)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error ReadingFile: %s %s", path, err.Error())
	}
	meta = new(FileMeta)
	err = json.Unmarshal(data, meta)
	return
}

//添加文件meta信息
func (node *RaftNode) CreateFileMeta(key string, meta *FileMeta) (err error) {
	if node.raft.State() != raft.Leader {
		return fmt.Errorf(ErrNotLeader)
	}
	if strings.Contains(key, "..") {
		return errors.New("invalid param")
	}

	c := &command{
		Op:   CommandCreateFileMeta,
		Key:  key,
		Meta: meta,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := node.raft.Apply(b, raftTimeout)
	return f.Error()
}

//更新文件meta信息
func (node *RaftNode) UpdateFileMeta(key string, meta *FileMeta) (err error) {
	if node.raft.State() != raft.Leader {
		return fmt.Errorf(ErrNotLeader)
	}
	if strings.Contains(key, "..") {
		return errors.New("invalid param")
	}

	c := &command{
		Op:   CommandUpdateFileMeta,
		Key:  key,
		Meta: meta,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := node.raft.Apply(b, raftTimeout)
	return f.Error()
}

//删除文件meta信息
func (node *RaftNode) DeleteFileMeta(key string) (err error) {
	if node.raft.State() != raft.Leader {
		return fmt.Errorf(ErrNotLeader)
	}
	if strings.Contains(key, "..") {
		return errors.New("invalid param")
	}

	c := &command{
		Op:   CommandDeleteFileMeta,
		Key:  key,
		Meta: nil,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := node.raft.Apply(b, raftTimeout)
	return f.Error()
}

//获取DataNode节点信息
func (node *RaftNode) GetDataNodeList() []*model.DataNode {
	node.mu.Lock()
	defer node.mu.Unlock()
	dataNodeList := make([]*model.DataNode, len(node.DataNodeSet))
	i := 0
	for _, v := range node.DataNodeSet {
		dataNodeList[i] = v
		i++
	}
	//fmt.Println(dataNodeList)
	return dataNodeList
}

//获取NameNode节点信息
func (node *RaftNode) GetNameNodeList() []*model.NameNode {
	node.mu.Lock()
	defer node.mu.Unlock()
	nameNodeList := make([]*model.NameNode, len(node.NameNodeSet))
	i := 0
	for _, v := range node.NameNodeSet {
		//fmt.Println("遍历节点", v.NodeID, v.HAddr)
		nameNodeList[i] = v
		i++
	}
	//fmt.Println(dataNodeList)
	return nameNodeList
}

//获取leaderNameNode节点的http地址
func (node *RaftNode) GetLeaderNameNode() *model.NameNode {
	_, nodeID := node.raft.LeaderWithID()

	//fmt.Println("leader节点为: ", nodeID, string(nodeID), node.NameNodeSet[string(nodeID)])

	node.mu.Lock()
	defer node.mu.Unlock()
	return node.NameNodeSet[string(nodeID)]
}

//添加DataNode节点，如果节点存在则，更新对应的节点数据
func (node *RaftNode) AddDataNode(dataNode *model.DataNode) (err error) {
	if node.raft.State() != raft.Leader {
		return fmt.Errorf(ErrNotLeader)
	}
	c := &command{
		Op:       CommandAddDataNode,
		DataNode: dataNode,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := node.raft.Apply(b, raftTimeout)
	return f.Error()
}

//添加NameNode节点的http服务地址
func (node *RaftNode) AddNameNodeHaddr(nameNode *model.NameNode) (err error) {
	if node.raft.State() != raft.Leader {
		return fmt.Errorf(ErrNotLeader)
	}
	c := &command{
		Op:       CommandAddNameNodeHaddr,
		NameNode: nameNode,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f := node.raft.Apply(b, raftTimeout)
	return f.Error()
}

//加入NameNode节点集群
func join(joinAddr, raftAddr, localHaddr, nodeID string) error {
	b, err := json.Marshal(model.ParamJoin{
		RaftAddr: raftAddr,
		HttpAddr: localHaddr,
		ID:       nodeID,
	})
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

func New(raftDir, metadir, raftAddr, joinAddr, localHaddr, nodeID string) (err error) {
	RaftNodeClient = &RaftNode{
		fileMutexs:  make(map[string]*sync.RWMutex),
		logger:      log.New(os.Stderr, "[NameNode] ", log.LstdFlags),
		RaftDir:     raftDir,
		MetaDir:     metadir,
		RaftBind:    raftAddr,
		DataNodeSet: make(map[string]*model.DataNode),
		NameNodeSet: make(map[string]*model.NameNode),
	}
	// RaftNodeClient.NameNodeSet[nodeID] = &model.NameNode{ //join haddr myself
	// 	NodeID: nodeID,
	// 	HAddr:  localHaddr,
	// }
	if err := RaftNodeClient.Open(joinAddr == "", nodeID); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
		return err
	}
	// If join was specified, make the join request.
	if joinAddr != "" {
		if err := join(joinAddr, raftAddr, localHaddr, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
			return err
		}
	} else { //join leader haddr info myself
		for {
			if RaftNodeClient.raft.State() == raft.Leader {
				err = RaftNodeClient.AddNameNodeHaddr(&model.NameNode{
					NodeID: nodeID,
					HAddr:  localHaddr,
				})
				if err != nil {
					log.Fatalf("join haddr myself err: {%s}", err.Error())
				}
				break
			}
			time.Sleep(time.Second)
		}
	}

	return nil
}

func (s *RaftNode) Open(enableSingle bool, localID string) error {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(s.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore

	boltDB, err := raftboltdb.New(raftboltdb.Options{
		Path: filepath.Join(s.RaftDir, "raft.db"),
	})
	if err != nil {
		return fmt.Errorf("new bbolt store: %s", err)
	}
	logStore = boltDB
	stableStore = boltDB

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, (*fsm)(s), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	s.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	return nil
}

func (s *RaftNode) Join(nodeID, addr, haddr string) error {
	s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.logger.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}

	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}

type fsm RaftNode

type command struct {
	Op       string          `json:"op,omitempty"`
	Key      string          `json:"key,omitempty"`
	Meta     *FileMeta       `json:"file_meta,omitempty"`
	DataNode *model.DataNode `json:"data_node,omitempty"`
	NameNode *model.NameNode `json:"name_node,omitempty"`
}

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case CommandCreateFileMeta:
		return f.applyCreateMeta(c.Key, c.Meta)
	case CommandUpdateFileMeta:
		return f.applyUpdateMeta(c.Key, c.Meta)
	case CommandDeleteFileMeta:
		return f.applyDeleteMeta(c.Key)
	case CommandAddDataNode:
		return f.applyAddDataNode(c.DataNode)
	case CommandAddNameNodeHaddr:
		return f.applyAddNameNodeHaddr(c.NameNode)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{}, nil
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	return nil
}

func (f *fsm) applyCreateMeta(key string, meta *FileMeta) interface{} {
	f.mu.Lock()
	fm, ok := f.fileMutexs[key]
	if !ok {
		f.fileMutexs[key] = &sync.RWMutex{}
		fm = f.fileMutexs[key]
	}
	f.mu.Unlock()
	fm.Lock()
	defer fm.Unlock()

	path := filepath.Join(f.MetaDir, key)
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("序列化meta信息失败")
	}
	os.WriteFile(path, data, 0644)
	return nil
}

func (f *fsm) applyUpdateMeta(key string, meta *FileMeta) interface{} {
	f.mu.Lock()
	fm, ok := f.fileMutexs[key]
	if !ok {
		f.fileMutexs[key] = &sync.RWMutex{}
		fm = f.fileMutexs[key]
	}
	f.mu.Unlock()
	fm.Lock()
	defer fm.Unlock()
	path := filepath.Join(f.MetaDir, key)
	_, err := os.Stat(path)
	if err != nil {
		logger.Logger.Error("文件key不存在", zap.Error(err))
		return errors.New("文件key不存在")
	}
	data, err := json.Marshal(meta)
	if err != nil {
		logger.Logger.Error("序列化meta信息失败", zap.Error(err))
		return errors.New("序列化meta信息失败")
	}
	os.WriteFile(path, data, 0644)
	return nil
}

func (f *fsm) applyDeleteMeta(key string) interface{} {
	f.mu.Lock()
	fm, ok := f.fileMutexs[key]
	if !ok {
		f.fileMutexs[key] = &sync.RWMutex{}
		fm = f.fileMutexs[key]
	}
	f.mu.Unlock()
	fm.Lock()
	defer fm.Unlock()

	path := filepath.Join(f.MetaDir, key)
	err := os.Remove(path)
	return err
}

func (f *fsm) applyAddDataNode(dataNode *model.DataNode) interface{} {
	f.mu.Lock()
	f.DataNodeSet[dataNode.URL] = dataNode
	f.mu.Unlock()
	log.Printf("new DataNode(host:%s name:%s nodedisk:%d free:%d) Join in", dataNode.URL, dataNode.NodeName, dataNode.NodeDiskSize, dataNode.NodeDiskAvailableSize)
	return nil
}

func (f *fsm) applyAddNameNodeHaddr(nameNode *model.NameNode) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	//fmt.Println("加入节点2:", nameNode.NodeID, nameNode.HAddr)
	f.NameNodeSet[nameNode.NodeID] = nameNode
	return nil
}

type fsmSnapshot struct{}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (f *fsmSnapshot) Release() {}
