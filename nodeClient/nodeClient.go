package nodeClient

import (
	"LDFS/model"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	ErrNotFoundUrl = errors.New("URL不存在或未初始化配置")
)

const (
	StoragePolicyEC   string = "EC"
	StoragePolicyCopy string = "cpoy"
)

//NameNodeHttpClient
type NameNodeHttpClient struct {
	getAllFileKeysUrl       string
	getFileMetaByFileKeyUrl string
	requestUploadFileUrl    string
	completeSampleUploadUrl string
	getDataNodeListInfoUrl  string
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
		getAllFileKeysUrl:       "/LDFS/getAllFileKeys",
		getFileMetaByFileKeyUrl: "/LDFS/getFileMetaByFileKey",
		requestUploadFileUrl:    "/LDFS/requestUploadFile",
		completeSampleUploadUrl: "/LDFS/completeSampleUpload",
		getDataNodeListInfoUrl:  "/LDFS/getDataNodeListInfo",
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

//获取简单上传文件的分块DataNode地址列表
func (nameNodeClient *NameNodeHttpClient) GetDataNodeListInfo(backendUrl string) (dataNodeList []*model.DataNode, err error) {
	res, err := http.Get(backendUrl + nameNodeClient.getDataNodeListInfoUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	resBytes, err := ioutil.ReadAll(res.Body)
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
	if resp.StatusCode != http.StatusOK {
		return errors.New("upload initiation failed with status " + resp.Status)
	}
	return
}

//获取文件完整FileMate信息
func (nameNodeClient *NameNodeHttpClient) GetFileMate(fileKey, backendUrl string) (meta *model.FileMetadata, err error) {
	//访问NameNode获取文件Meta信息
	res, err := http.Get(backendUrl + nameNodeClient.getFileMetaByFileKeyUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// 检查响应状态码
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("file download failed with status " + res.Status)
	}
	fileMetaBytes, err := ioutil.ReadAll(res.Body)
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
		return nil, errors.New("list parts failed with status " + resp.Status)
	}

	// 解析响应体
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
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
	req, err := http.NewRequest("POST", URL, nil)
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
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
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

//恢复文件数据
func (dataNodeClient *DataNodeHttpClient) RecoverShard(opt *model.RecoverShardParam, backend string) (err error) {
	URL := backend + dataNodeClient.recoverShardUrl

	requestBodyBytes, err := json.Marshal(opt)
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
	if resp.StatusCode != http.StatusOK {
		return errors.New("get fileKeys list failed with status " + resp.Status)
	}

	return nil
}
