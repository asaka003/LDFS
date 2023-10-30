package storagesdk

import (
	"LDFS/fileNode/util"
	"LDFS/model"
	"LDFS/nodeClient"
	"bytes"
	"encoding/json"
	"io"
	"math"
	"os"
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

//初始化分块上传文件
func (cli *ObjectClient) InitiateMultipartUpload(fileKey string, fileSize int64) (UploadID string, err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	//请求NameNode初始化上传，获取文件上传节点
	UploadID, err = NameNodeClient.InitMultiUpload(fileKey, "", backend, fileSize)
	if err != nil {
		return
	}
	return
}

//上传分块
func (cli *ObjectClient) UploadPart(fileKey, uploadID string, partNumber int, r io.Reader) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}

	err = NameNodeClient.UploadMultiPart(fileKey, uploadID, backend, partNumber, r)
	return
}

//查询已上传分块
func (cli *ObjectClient) ListParts(fileKey, uploadID string) (parts []*model.Object, err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}

	parts, err = NameNodeClient.ListParts(fileKey, uploadID, backend)
	return
}

//完成分块上传
func (cli *ObjectClient) CompleteMultipartUpload(fileKey, uploadID string, opt *model.CompleteMultipartUploadOptions) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	err = NameNodeClient.CompeleteMultiUpload(fileKey, "", uploadID, backend, opt)
	return
}

//终止分块上传
func (cli *ObjectClient) AbortMultipartUpload(fileKey, uploadID string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	err = NameNodeClient.AbortMultipartUpload(fileKey, uploadID, backend)
	return
}

//副本冗余模式Client
type ReplicasClient struct {
	ObjectClient
}

//初始化副本冗余模式Client
func NewReplicasClient(nameNodeUrls []string) (client StorageClient) {
	NameNodeUrls = nameNodeUrls
	InitCluster(nameNodeUrls)
	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	return &ReplicasClient{}
}

//下载文件
func (cli *ReplicasClient) DownloadFile(fileKey string, destPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	fileMeta, err := NameNodeClient.GetFileMate(fileKey, backend)
	if err != nil {
		return
	}
	shards := fileMeta.Shards
	// 打开文件以附加模式打开，如果文件不存在则创建
	file, err := os.OpenFile(destPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	for _, shard := range shards {
		err = DataNodeClient.ReplicasDownloadShard(shard.Hash, shard.NodeURLs[0], file)
		if err != nil {
			return err
		}
	}

	return
}

//简单上传文件
func (cli *ReplicasClient) SimpleUploadFile(fileKey string, srcPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
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

	fileMeta, err := NameNodeClient.RequestUploadFile(fileKey, backend, StoragePolicyCopy, fileInfo.Size(), Copy_BlockSize)
	if err != nil {
		return
	}

	//上传文件块
	buf := make([]byte, Copy_BlockSize)
	buffer := bytes.NewBuffer(buf)
	for _, shard := range fileMeta.Shards {
		// 重置 buffer，以便重新使用
		buffer.Reset()
		// 分块读取文件数据到 buffer
		_, err = io.CopyN(buffer, file, int64(Copy_BlockSize))
		if err != nil && err != io.EOF {
			return
		}
		shard.Hash = util.BytesHash(buffer.Bytes())
		var shardJson []byte
		shardJson, err = json.Marshal(shard)
		if err != nil {
			return
		}
		err = DataNodeClient.ReplicasUploadShard(shard.Hash, string(shardJson), buffer, shard.NodeURLs[0], 0)
		if err != nil {
			return
		}
	}

	return
}

//EC模式-条形布局Client

//EC模式-连续布局Client
type ContiguousLayoutClient struct {
}

//初始化SDK列表,加载DataNode地址列表,目前只支持http协议(连续布局策略)
func NewContiguousLayoutClient(nameNodeUrls []string) (client StorageClient) {
	//DataNodeUrls = dataNodeUrls
	NameNodeUrls = nameNodeUrls
	//InitReverseProxy(DataNodeUrls)
	InitCluster(nameNodeUrls)
	NameNodeClient = nodeClient.GetNameNodeHttpClient()
	DataNodeClient = nodeClient.GetDataNodeHttpClient()
	return &ContiguousLayoutClient{}
}

//下载文件
func (cli *ContiguousLayoutClient) DownloadFile(fileKey string, destPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	fileMeta, err := NameNodeClient.GetFileMate(fileKey, backend)
	if err != nil {
		return err
	}

	// //生成随机目录
	// randomInPath := filepath.Join(EC_InputDir, generateRandomString(6))
	// err = os.MkdirAll(randomInPath, os.ModeDir)
	// if err != nil {
	// 	return
	// }
	// defer func() { //移除临时编码文件
	// 	os.RemoveAll(randomInPath)
	// }()
	// filePaths := []string{}
	// for _, shard := range fileMeta.Shards {
	// 	fpath := filepath.Join(randomInPath, shard.Hash)
	// 	filePaths = append(filePaths, fpath)
	// 	DataNodeClient.downloadShard(shard.Hash, shard.NodeURLs[0], fpath)
	// }

	blockNum := len(fileMeta.Shards)
	buffers := make([]*bytes.Buffer, fileMeta.DataShards+fileMeta.ParityShards)
	for i := range buffers {
		buf := make([]byte, EC_ShardSize)
		buffers[i] = bytes.NewBuffer(buf)
	}
	//下载文件Shard
	for blockId := 0; blockId < blockNum; blockId++ {
		length := fileMeta.DataShards + fileMeta.ParityShards
		for shardId := 0; shardId < length; shardId++ {
			err = DataNodeClient.ECDownloadShard(fileMeta.Shards[blockId].PartHashs[shardId], fileMeta.Shards[blockId].NodeURLs[shardId], buffers[shardId])
			if err != nil {
				return err
			}
		}
	}

	err = ReconstructBuffer(buffers, destPath, fileMeta.DataShards, fileMeta.ParityShards, EC_BlockSize)
	return
}

//简单上传文件
func (cli *ContiguousLayoutClient) SimpleUploadFile(fileKey string, srcPath string) (err error) {
	//一致性hash负载均衡获取要请求的NameNode地址
	backend, err := nameNodeCluster.Consistent.Get(fileKey)
	if err != nil {
		return
	}
	srcFile, err := os.OpenFile(srcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	fileInfo, err := srcFile.Stat()
	if err != nil {
		return
	}
	//请求NameNode初始化上传，获取文件上传节点
	fileMeta, err := NameNodeClient.RequestUploadFile(fileKey, backend, StoragePolicyEC, fileInfo.Size(), EC_BlockSize)
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
		io.CopyN(blockBuf, srcFile, EC_BlockSize)
		//编码文件
		shardBuffs, err := EncodeBuffer(blockBuf, fileMeta.DataShards, fileMeta.ParityShards)
		if err != nil {
			return err
		}
		length := len(shardBuffs)
		for i := 0; i < length; i++ {
			//计算shard的Hash值
			Hash := util.BytesHash(shardBuffs[i].Bytes())
			fileMeta.Shards[blockId].PartHashs[i] = Hash
			err = DataNodeClient.ECUploadShardBytes(Hash, blockBuf, fileMeta.Shards[blockId].NodeURLs[i])
			if err != nil {
				return err
			}
		}
	}

	//将保存文件DataNode存储列表信息
	err = NameNodeClient.SendSampleUploadInfo(fileKey, backend, fileMeta.Shards)
	if err != nil {
		return
	}
	return
}
