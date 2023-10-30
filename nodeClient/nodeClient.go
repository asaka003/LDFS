package nodeClient

import (
	"LDFS/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	getSampleUploadListUrl  string
	sendSampleUploadInfoUrl string
	getFileMateUrl          string

	initMultiUploadUrl      string
	uploadMultiPartUrl      string
	compeleteMultiUploadUrl string
	abortMultipartUploadUrl string
	listPartsUrl            string

	requestUploadFileUrl string
}

type DataNodeHttpClient struct {
	ECuploadShardUrl         string
	ECdownloadShardUrl       string
	replicasUploadShardUrl   string
	replicasDownloadShardUrl string
}

func GetNameNodeHttpClient() *NameNodeHttpClient {
	return &NameNodeHttpClient{
		getSampleUploadListUrl:  "/LDFS/getSampleUploadList",
		sendSampleUploadInfoUrl: "/LDFS/sendSampleUploadInfo",
		getFileMateUrl:          "/LDFS/getFileMeta",

		initMultiUploadUrl:      "/LDFS/initMultiUpload",
		uploadMultiPartUrl:      "/LDFS/uploadMultiPart",
		compeleteMultiUploadUrl: "/LDFS/compeleteMultiUpload",
		abortMultipartUploadUrl: "/LDFS/abortMultipartUpload",
		listPartsUrl:            "/LDFS/listParts",
		requestUploadFileUrl:    "/LDFS/requestUploadFile",
	}
}

func GetDataNodeHttpClient() *DataNodeHttpClient {
	return &DataNodeHttpClient{
		ECuploadShardUrl:         "/LDFS/ECuploadShard",
		ECdownloadShardUrl:       "/LDFS/ECdownloadShard",
		replicasUploadShardUrl:   "/LDFS/replicasUploadShard",
		replicasDownloadShardUrl: "/LDFS/replicasDownloadShard",
	}
}

//获取简单上传文件的分块DataNode地址列表
func (nameNodeClient *NameNodeHttpClient) GetSampleUploadList(backendUrl string) (urlList *model.SampleUploadList, err error) {
	res, err := http.Get(backendUrl + nameNodeClient.getSampleUploadListUrl)
	if err != nil {
		return
	}
	defer res.Body.Close()
	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	urlList = new(model.SampleUploadList)
	err = json.Unmarshal(resBytes, urlList)
	if err != nil {
		return
	}
	return
}

//发送简单上传分块信息
func (nameNodeClient *NameNodeHttpClient) SendSampleUploadInfo(fileKey, backendUrl string, list []*model.Shard) (err error) {
	URL := backendUrl + nameNodeClient.sendSampleUploadInfoUrl

	// 创建请求体
	requestBody := model.SampleUploadInfo{FileKey: fileKey, Shards: list}
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
	res, err := http.Get(backendUrl + nameNodeClient.getFileMateUrl)
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

//初始化分块上传文件()
func (nameNodeClient *NameNodeHttpClient) InitMultiUpload(fileKey, fileHash, backendUrl string, fileSize int64) (UploadID string, err error) {
	return
}

//上传文件分块
func (nameNodeClient *NameNodeHttpClient) UploadMultiPart(fileKey, uploadID, backendUrl string, partNumber int, r io.Reader) (err error) {
	URL := backendUrl + nameNodeClient.uploadMultiPartUrl

	// 创建一个管道
	pr, pw := io.Pipe()

	// 创建一个multipart.Writer，它将写入管道
	multipartWriter := multipart.NewWriter(pw)

	// 在一个单独的goroutine中处理文件的写入
	go func() {
		defer pw.Close()
		defer multipartWriter.Close()

		// 添加"FileKey"字段
		if err := multipartWriter.WriteField("FileKey", fileKey); err != nil {
			fmt.Println("Error adding field 'FileKey':", err)
			return
		}

		// 添加"UploadID"字段
		if err := multipartWriter.WriteField("UploadID", uploadID); err != nil {
			fmt.Println("Error adding field 'UploadID':", err)
			return
		}

		// 添加"ChunkIndex"字段
		if err := multipartWriter.WriteField("ChunkIndex", fmt.Sprintf("%d", partNumber)); err != nil {
			fmt.Println("Error adding field 'ChunkIndex':", err)
			return
		}

		// 添加文件数据
		part, err := multipartWriter.CreateFormFile("file", fileKey)
		if err != nil {
			fmt.Println("Error adding field 'file':", err)
			return
		}

		if _, err := io.Copy(part, r); err != nil {
			fmt.Println("Error copying file to part:", err)
			return
		}
	}()

	// 构建请求
	req, err := http.NewRequest("POST", URL, pr)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return errors.New("upload part failed with status " + resp.Status)
	}

	return nil
}

//完成分块上传
func (nameNodeClient *NameNodeHttpClient) CompeleteMultiUpload(fileKey, fileHash, uploadID, backend string, opt *model.CompleteMultipartUploadOptions) (err error) {
	URL := backend + nameNodeClient.compeleteMultiUploadUrl

	// 创建请求体
	requestBody := model.CompleteMultipartParam{UploadID: uploadID, FileKey: fileKey, FileHash: fileHash, Opt: opt}
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
		return errors.New("complete upload failed with status " + resp.Status)
	}
	return
}

//终止分块上传文件
func (nameNodeClient *NameNodeHttpClient) AbortMultipartUpload(fileKey, uploadID, backend string) (err error) {
	URL := backend + nameNodeClient.abortMultipartUploadUrl

	// 创建请求体
	requestBody := model.AbortMultipartUploadParam{UploadID: uploadID, FileKey: fileKey}
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
		return errors.New("complete upload failed with status " + resp.Status)
	}
	return
}

//查询已上传分块
func (nameNodeClient *NameNodeHttpClient) ListParts(fileKey, uploadID, backend string) (parts []*model.Object, err error) {
	URL := backend + nameNodeClient.listPartsUrl

	// 创建请求体
	requestBody := model.ListPartsParam{UploadID: uploadID, FileKey: fileKey}
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
	var response = new(model.ListPartsResponse)
	err = json.Unmarshal(responseBodyBytes, response)
	if err != nil {
		return
	}

	// 返回 partNumbers
	return response.Parts, nil
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
