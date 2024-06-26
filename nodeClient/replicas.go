package nodeClient

import (
	"LDFS/model"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

/*副本冗余模式*/

//上传文件数据块
func (dataNodeClient *DataNodeHttpClient) ReplicasUploadShard(shardHash, blockJson string, data *bytes.Buffer, backend string, copyNo int64) (err error) {
	// 创建一个管道
	pr, pw := io.Pipe()
	defer pr.Close()
	// 创建一个multipart.Writer，它将写入管道
	multipartWriter := multipart.NewWriter(pw)

	// 在一个单独的goroutine中处理文件的写入
	go func() {
		defer pw.Close()
		defer multipartWriter.Close() //

		// 添加"FileKey"字段
		if err = multipartWriter.WriteField("hash", shardHash); err != nil {
			return
		}

		//添加shards信息
		if err = multipartWriter.WriteField("blockJson", blockJson); err != nil {
			return
		}

		//添加copyNo信息
		if err = multipartWriter.WriteField("copyNo", strconv.FormatInt(copyNo, 10)); err != nil {
			return
		}

		// 添加文件数据
		var part io.Writer
		part, err = multipartWriter.CreateFormFile("file", shardHash)
		if err != nil {
			return
		}

		_, err = io.Copy(part, data)

		multipartWriter.Close() //及时关闭，给报文添加结束行
	}()

	// 构建请求
	req, err := http.NewRequest("POST", backend+dataNodeClient.replicasUploadShardUrl, pr)
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
		return errors.New("upload failed with status " + resp.Status)
	}
	return nil
}

//下载文件数据块
func (dataNodeClient *DataNodeHttpClient) ReplicasDownloadShard(fileKey string, blockID int, shardHash, backend string, shard *model.Shard, des io.Writer) (err error) {
	URL := backend + dataNodeClient.replicasDownloadShardUrl
	// 创建请求体
	requestBody := model.ParamReplicasDownloadShard{
		FileKey: fileKey,
		BlockId: blockID,
		Shard:   shard,
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
		return errors.New("list parts failed with status " + resp.Status)
	}

	io.Copy(des, resp.Body)
	return nil
}
