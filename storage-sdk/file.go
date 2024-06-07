package storagesdk

import (
	"LDFS/dataNode/util"
	"LDFS/model"
	"LDFS/nodeClient"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

//不同存储策略抽象接口
type StorageClient interface {
	DownloadFile(fileKey string, destPath string) (err error)
	SimpleUploadFile(fileKey string, srcPath string) (err error)
	//SimpleUploadFile(fileKey string, r io.Reader) (err error)
}

//父类存储客户端
type ObjectClient struct {
}

//副本冗余模式Client
type ReplicasClient struct {
	ObjectClient
}

//初始化副本冗余模式Client
func NewReplicasClient(nameNodeLeaderUrls, nameNodeFollowerUrls []string) (client StorageClient) {
	NameNodeLeaderUrls = nameNodeLeaderUrls
	NameNodeFollowerUrls = nameNodeFollowerUrls
	InitCluster(NameNodeLeaderUrls, nameNodeFollowerUrls)
	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	return &ReplicasClient{}
}

//下载文件
func (cli *ReplicasClient) DownloadFile(fileKey string, destPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeFollowerCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	fileMeta, err := NameNodeClient.GetFileMate(fileKey, backend)
	if err != nil {
		return
	}
	blocks := fileMeta.Blocks
	// 打开文件以附加模式打开，如果文件不存在则创建
	file, err := os.OpenFile(destPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	wg := sync.WaitGroup{}
	for _, block := range blocks {
		ok := false
		for _, shard := range block.Shards {
			shard.Hash = block.Hash
			err = DataNodeClient.ReplicasDownloadShard(fileKey, block.BlockId, block.Hash, shard.NodeURL, shard, file)
			if err != nil {
				log.Printf("file broken:%s", err.Error())
				// wg.Add(1)
				// go func() {
				// 	DataNodeClient.RecoverShard(&model.RecoverShardParam{ //请求DataNode修复文件
				// 		Block:         block,
				// 		ShardId:       shard.ShardID,
				// 		StoragePolicy: nodeClient.StoragePolicyCopy,
				// 	}, shard.NodeURL)
				// 	wg.Done()
				// }()
				// continue
			} else {
				ok = true
				break
			}
		}
		if !ok {
			return errors.New("文件损坏")
		}
	}
	wg.Wait()
	return
}

//简单上传文件
func (cli *ReplicasClient) SimpleUploadFile(fileKey string, srcPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeLeaderCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	file, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer file.Close()
	fileInfo, err := os.Stat(srcPath)
	if err != nil {
		return
	}
	fileMeta, err := NameNodeClient.RequestUploadFile(fileKey, backend, nodeClient.StoragePolicyCopy, fileInfo.Size(), Copy_BlockSize)
	if err != nil {
		//变更leader
		var nameNode *model.NameNode
		for _, url := range NameNodeFollowerUrls {
			nameNode, err = NameNodeClient.GetNameNodeLeaderInfo(url)
			if err != nil {
				log.Printf("get nameNode Err")
			} else {
				break
			}
		}
		AddNameNodeLeader("http://" + nameNode.HAddr)
		RemoveNameNodeLeader(backend)
		backend, err = nameNodeLeaderCluster.Consistent.Get(fileKey)
		if err != nil {
			return err
		}
		fileMeta, err = NameNodeClient.RequestUploadFile(fileKey, backend, nodeClient.StoragePolicyCopy, fileInfo.Size(), Copy_BlockSize)
		if err != nil {
			return err
		}
	}
	//上传文件块
	buf := make([]byte, 0) //置0,利用buffer动态扩容机制
	buffer := bytes.NewBuffer(buf)
	for _, block := range fileMeta.Blocks {
		// 重置 buffer，以便重新使用
		buffer.Reset()
		// 分块读取文件数据到 buffer
		_, err = io.CopyN(buffer, file, int64(Copy_BlockSize))
		if err != nil && err != io.EOF {
			return
		}
		block.Hash = util.BytesHash(buffer.Bytes())
		var blockJson []byte
		blockJson, err = json.Marshal(block)
		if err != nil {
			return
		}
		err = DataNodeClient.ReplicasUploadShard(block.Hash, string(blockJson), buffer, block.Shards[0].NodeURL, 0)
		if err != nil {
			return err
		}
		//更新文件meta信息
		err = NameNodeClient.UpdateFileMeta(backend, fileMeta)
		if err != nil {
			if err.Error() == nodeClient.ErrNotLeader {
				var nameNode *model.NameNode
				for _, url := range NameNodeFollowerUrls {
					nameNode, err = NameNodeClient.GetNameNodeLeaderInfo(url)
					if err != nil {
						log.Printf("get nameNode Err")
					} else {
						break
					}
				}
				AddNameNodeLeader(nameNode.HAddr)
				RemoveNameNodeLeader(backend)
				backend, err = nameNodeLeaderCluster.Consistent.Get(fileKey)
				if err != nil {
					return err
				}
				err = NameNodeClient.UpdateFileMeta(nameNode.HAddr, fileMeta) //再次尝试更新Meta信息
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return
}

//EC模式-Client
type ECClient struct {
}

//初始化SDK列表,目前只支持http协议
func NewECClient(nameNodeLeaderUrls, nameNodeFollowerUrls []string) (client StorageClient) {
	NameNodeLeaderUrls = nameNodeLeaderUrls
	NameNodeFollowerUrls = nameNodeFollowerUrls
	InitCluster(NameNodeLeaderUrls, nameNodeFollowerUrls)
	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	return &ECClient{}
}

//下载文件
func (cli *ECClient) DownloadFile(fileKey string, destPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeFollowerCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	fileMeta, err := NameNodeClient.GetFileMate(fileKey, backend)
	if err != nil {
		return err
	}

	blockNum := len(fileMeta.Blocks)
	buffers := make([]*bytes.Buffer, fileMeta.DataShards+fileMeta.ParityShards)
	for i := range buffers {
		buffers[i] = new(bytes.Buffer)
	}
	//下载文件Shard
	wg := sync.WaitGroup{}
	for blockId := 0; blockId < blockNum; blockId++ {
		length := fileMeta.DataShards + fileMeta.ParityShards
		for shardId := 0; shardId < length; shardId++ {
			err = DataNodeClient.ECDownloadShard(fileMeta.Blocks[blockId].Shards[shardId].Hash, fileMeta.Blocks[blockId].Shards[shardId].NodeURL, buffers[shardId])
			if err != nil {
				log.Printf("file broken:%s", err.Error())
				return err
			}
			// if err != nil {
			// 	buffers[shardId] = nil
			// 	wg.Add(1)
			// 	go func() {
			// 		DataNodeClient.RecoverShard(&model.RecoverShardParam{ //请求DataNode修复文件
			// 			Block:          fileMeta.Blocks[blockId],
			// 			ShardId:        int64(shardId),
			// 			DataShardNum:   fileMeta.DataShards,
			// 			ParityShardNum: fileMeta.ParityShards,
			// 			StoragePolicy:  nodeClient.StoragePolicyEC,
			// 		}, fileMeta.Blocks[blockId].Shards[shardId].NodeURL)
			// 		wg.Done()
			// 	}()
			// 	return err
			// }
		}
	}

	err = ReconstructBuffer(buffers, destPath, fileMeta.DataShards, fileMeta.ParityShards, EC_BlockSize)
	wg.Wait()
	return
}

/*
	纠删码模式下的 简单上传文件
	此方法会将用户给定的源文件地址进行block_size大小的切分，
	随后对每一个文件块进行数据块和校验块计算，并分发指定的DataNode节点上
*/
func (cli *ECClient) SimpleUploadFile(fileKey string, srcPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeLeaderCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return
	}
	//请求NameNode初始化上传，获取文件上传节点
	fileMeta, err := NameNodeClient.RequestUploadFile(fileKey, backend, nodeClient.StoragePolicyEC, fileInfo.Size(), EC_BlockSize)
	if err != nil {
		return
	}

	BlockNum := int(math.Ceil(float64(fileInfo.Size()) / float64(EC_BlockSize)))
	for blockId := 0; blockId < BlockNum; blockId++ {

	}

	blockData := make([]byte, EC_BlockSize)
	blockBuf := bytes.NewBuffer(blockData)
	for blockId := 0; blockId < BlockNum; blockId++ {
		blockBuf.Reset()
		//读取block数据
		_, err := io.CopyN(blockBuf, srcFile, EC_BlockSize)
		if err != io.EOF && err != nil {
			return err
		}
		//编码文件
		shardBuffs, err := EncodeBuffer(blockBuf, fileMeta.DataShards, fileMeta.ParityShards)
		if err != nil {
			return err
		}
		length := len(shardBuffs)
		for i := 0; i < length; i++ {
			//计算shard的Hash值
			Hash := util.BytesHash(shardBuffs[i].Bytes())
			fileMeta.Blocks[blockId].Shards[i].Hash = Hash
			err = DataNodeClient.ECUploadShardBytes(Hash, shardBuffs[i], fileMeta.Blocks[blockId].Shards[i].NodeURL)
			if err != nil {
				return err
			}
		}
	}

	//完成文件EC上传，并更新filemeta信息
	err = NameNodeClient.CompleteSampleUpload(fileKey, backend, fileMeta)
	if err != nil {
		if err.Error() == nodeClient.ErrNotLeader {
			var nameNode *model.NameNode
			for _, url := range NameNodeFollowerUrls {
				nameNode, err = NameNodeClient.GetNameNodeLeaderInfo(url)
				if err != nil {
					log.Printf("get nameNode Err")
				} else {
					break
				}
			}
			AddNameNodeLeader(nameNode.HAddr)
			RemoveNameNodeLeader(backend)
			backend, err = nameNodeLeaderCluster.Consistent.Get(fileKey)
			if err != nil {
				return err
			}
			err = NameNodeClient.UpdateFileMeta(nameNode.HAddr, fileMeta) //再次尝试更新Meta信息
			if err != nil {
				return err
			}
		}
		return
	}
	return
}
