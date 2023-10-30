package util

import (
	"LDFS/fileNode/config"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrCopyFile  = errors.New("副本冗余复制文件失败")
	ErrCopyShard = errors.New("纠删码数据块冗余传输失败")
)

//流式处理文件发送
func sendCopyFileToNode(nodeURL, filePath, copyType string) error {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个管道
	pipeReader, pipeWriter := io.Pipe()

	// 创建一个multipart.Writer，它将写入管道
	multipartWriter := multipart.NewWriter(pipeWriter)

	// 在一个单独的goroutine中处理文件的写入
	go func() {
		defer pipeWriter.Close()
		defer multipartWriter.Close()

		// 添加"type"字段
		if err := multipartWriter.WriteField("type", copyType); err != nil {
			fmt.Println("Error adding field 'type':", err)
			return
		}

		// 添加文件数据
		part, err := multipartWriter.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			fmt.Println("Error adding field 'file':", err)
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			fmt.Println("Error copying file to part:", err)
			return
		}
	}()

	// 构建请求
	req, err := http.NewRequest("POST", nodeURL+"/copyFileData", pipeReader)
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

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + resp.Status)
	}
	return nil
}

//副本冗余模式处理文件上传
func MultiNodeCopyFile(urls []string, localUrl, filePath string) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	for _, nodeURL := range urls {
		wg.Add(1)
		go func(nodeURL string) {
			var err error
			if nodeURL == config.LocalFileNodeUrl {
				err = nil
			} else {
				err = sendCopyFileToNode(nodeURL, filePath, "file")
			}
			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			wg.Done()
		}(nodeURL)
	}
	wg.Wait()

	if successCount > len(urls)/2 {
		return nil
	} else {
		return ErrCopyFile
	}
}
