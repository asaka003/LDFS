package nodeClient

import (
	"LDFS/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrNotFoundUrl = errors.New("URL不存在或未初始化配置")
	ErrRequestErr  = errors.New("请求出错")
	ErrServerErr   = errors.New("服务端错误")
)

const (
	StoragePolicyEC   string = "EC"
	StoragePolicyCopy string = "copy"
)

//NameNodeHttpClient
type NameNodeHttpClient struct {
	getAllFileKeysUrl        string
	getFileMetaByFileKeyUrl  string
	requestUploadFileUrl     string
	completeSampleUploadUrl  string
	getDataNodeListInfoUrl   string
	updateFileMetaUrl        string
	getNameNodeLeaderInfoUrl string
	getNameNodeListInfoUrl   string
	joinNameNodeUrl          string
	joinDataNodeUrl          string
}

type DataNodeHttpClient struct {
	ECuploadShardUrl         string
	ECdownloadShardUrl       string
	replicasUploadShardUrl   string
	replicasDownloadShardUrl string
	recoverShardUrl          string
	getStorageInfoUrl        string
}

func GetNameNodeHttpClient() *NameNodeHttpClient {
	return &NameNodeHttpClient{
		getAllFileKeysUrl:        "/LDFS/getAllFileKeys",
		getFileMetaByFileKeyUrl:  "/LDFS/getFileMetaByFileKey",
		requestUploadFileUrl:     "/LDFS/requestUploadFile",
		completeSampleUploadUrl:  "/LDFS/completeSampleUpload",
		getDataNodeListInfoUrl:   "/LDFS/getDataNodeListInfo",
		updateFileMetaUrl:        "/LDFS/updateFileMeta",
		getNameNodeLeaderInfoUrl: "/LDFS/getNameNodeLeaderInfo",
		getNameNodeListInfoUrl:   "/LDFS/getNameNodeListInfo",
		joinNameNodeUrl:          "/LDFS/join",
		joinDataNodeUrl:          "/LDFS/joinDataNode",
	}
}

func GetDataNodeHttpClient() *DataNodeHttpClient {
	return &DataNodeHttpClient{
		ECuploadShardUrl:         "/LDFS/ECuploadShard",
		ECdownloadShardUrl:       "/LDFS/ECdownloadShard",
		replicasUploadShardUrl:   "/LDFS/replicasUploadShard",
		replicasDownloadShardUrl: "/LDFS/replicasDownloadShard",
		recoverShardUrl:          "/LDFS/recoverShard",
		getStorageInfoUrl:        "/LDFS/getStorageInfo",
	}
}

const ErrNotLeader = "err not leader"

//获取简单上传文件的分块DataNode地址列表
func (nameNodeClient *NameNodeHttpClient) GetDataNodeListInfo(backendUrl string) (dataNodeList []*model.DataNode, err error) {
	res, err := http.Get(backendUrl + nameNodeClient.getDataNodeListInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Get DataNode List Info failed with status " + res.Status)
		return
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	dataNodeList = make([]*model.DataNode, 0)
	err = json.Unmarshal(resBytes, &dataNodeList)
	if err != nil {
		return
	}
	return
}

//获取所有NameNode地址列表
func (nameNodeClient *NameNodeHttpClient) GetNameNodeListInfo(backend string) (nameNodeList []*model.NameNode, err error) {
	res, err := http.Get(backend + nameNodeClient.getNameNodeListInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Get DataNode List Info failed with status " + res.Status)
		return
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	nameNodeList = make([]*model.NameNode, 0)
	err = json.Unmarshal(resBytes, &nameNodeList)
	if err != nil {
		return
	}
	return
}

//加入DataNode节点
func (nameNodeClient *NameNodeHttpClient) JoinDataNode(leaderNameNodeAddr string, dataNode *model.DataNode) error {
	b, err := json.Marshal(model.ParamJoinDataNode{
		DataNodeInfo: dataNode,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s%s", leaderNameNodeAddr, nameNodeClient.joinDataNodeUrl), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

//发送简单上传分块信息
func (nameNodeClient *NameNodeHttpClient) CompleteSampleUpload(fileKey, backendUrl string) (err error) {
	URL := backendUrl + nameNodeClient.completeSampleUploadUrl

	// 创建请求体
	requestBody := model.CompleteSampleUploadParams{FileKey: fileKey}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode == http.StatusForbidden {
		return errors.New(ErrNotLeader)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Complete Sample Upload failed with status " + resp.Status)
	}
	return
}

//获取文件完整FileMate信息
func (nameNodeClient *NameNodeHttpClient) GetFileMate(fileKey, backendUrl string) (meta *model.FileMetadata, err error) {
	//访问NameNode获取文件Meta信息
	res, err := http.Get(backendUrl + nameNodeClient.getFileMetaByFileKeyUrl + "/" + fileKey)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Get FileMate failed with status " + res.Status)
	}
	fileMetaBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	meta = new(model.FileMetadata)
	err = json.Unmarshal(fileMetaBytes, meta)
	if err != nil {
		return
	}
	return
}

//请求上传文件(目前支持副本冗余模式)
func (nameNodeClient *NameNodeHttpClient) RequestUploadFile(fileKey, backend, storagePolicy string, fileSize, blockSize int64) (fileMeta *model.FileMetadata, err error) {
	URL := backend + nameNodeClient.requestUploadFileUrl

	// 创建请求体
	requestBody := model.RequestUploadFileParams{
		FileKey:       fileKey,
		FileSize:      fileSize,
		BlockSize:     blockSize,
		StoragePolicy: storagePolicy,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, ErrServerErr
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return nil, errors.New(ErrNotLeader)
	}
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request UploadFile failed with status " + resp.Status)
	}
	// 解析响应体
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var response = new(model.RequestUploadFileResponse)
	err = json.Unmarshal(responseBodyBytes, response)
	if err != nil {
		return
	}

	// 返回 各个block的存储列表
	return response.FileMeta, nil
}

//获取所有文件信息列表
func (nameNodeClient *NameNodeHttpClient) GetAllFileKeys(backend string) (fileList []*model.FileInfo, err error) {
	URL := backend + nameNodeClient.getAllFileKeysUrl

	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get fileKeys list failed with status " + resp.Status)
	}

	// 解析响应体
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result = make([]*model.FileInfo, 0)
	err = json.Unmarshal(responseBodyBytes, &result)
	if err != nil {
		return
	}

	return result, nil
}

//更新文件meta信息
func (nameNodeClient *NameNodeHttpClient) UpdateFileMeta(backend string, fileMeta *model.FileMetadata) (err error) {
	URL := backend + nameNodeClient.updateFileMetaUrl
	requestBodyBytes, err := json.Marshal(fileMeta)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode == http.StatusForbidden {
		return errors.New(ErrNotLeader)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Request UploadFile failed with status " + resp.Status)
	}
	return nil
}

//获取NameNode  Leaderhttp服务地址
func (nameNodeClient *NameNodeHttpClient) GetNameNodeLeaderInfo(backend string) (nameNode *model.NameNode, err error) {
	//访问NameNode获取文件Meta信息
	res, err := http.Get(backend + nameNodeClient.getNameNodeLeaderInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Get FileMate failed with status " + res.Status)
	}
	fileMetaBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	nameNode = new(model.NameNode)
	err = json.Unmarshal(fileMetaBytes, nameNode)
	if err != nil {
		return
	}
	return
}

//恢复文件数据
// func (dataNodeClient *DataNodeHttpClient) RecoverShard(opt *model.RecoverShardParam, backend string) (err error) {
// 	URL := backend + dataNodeClient.recoverShardUrl

// 	requestBodyBytes, err := json.Marshal(opt)
// 	if err != nil {
// 		return
// 	}
// 	// 创建 HTTP 请求
// 	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBodyBytes))
// 	if err != nil {
// 		return
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// 发送 HTTP 请求
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// 检查响应状态码
// 	if resp.StatusCode != http.StatusOK {
// 		return errors.New("Recover Shard failed with status " + resp.Status)
// 	}

// 	return nil
// }

//获取DataNode存储信息
func (dataNodeClient *DataNodeHttpClient) GetStorageInfo(backend string) (dataNodeInfo model.DataNode, err error) {
	res, err := http.Get(backend + dataNodeClient.getStorageInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Get Storage Info failed with status " + res.Status)
		return
	}
	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	dataNodeInfo = model.DataNode{}
	err = json.Unmarshal(resBytes, &dataNodeInfo)
	if err != nil {
		return
	}
	dataNodeInfo.URL = backend
	return
}
