package util

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

//合并分块上传文件
func MergeMultiFile(tmpMultiDir, newFilePath string) error {
	files, err := os.ReadDir(tmpMultiDir)
	if err != nil {
		return err
	}
	// 对文件信息按文件名排序
	sort.Slice(files, func(i, j int) bool {
		numI, _ := strconv.Atoi(files[i].Name())
		numJ, _ := strconv.Atoi(files[j].Name())
		return numI < numJ
	})

	// 创建新文件
	newFile, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	// 依次读取每个文件的内容，并将其写入到新文件中
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		// 打开分块文件
		filePath := filepath.Join(tmpMultiDir, file.Name())
		partFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		// 将文件内容写入到新文件中
		_, err = io.Copy(newFile, partFile)
		if err != nil {
			return err
		}
		partFile.Close()
	}
	err = os.RemoveAll(tmpMultiDir)
	return err
}
