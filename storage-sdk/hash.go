package storagesdk

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

//计算文件hash值
func getFileHash(path string) (Hash string, err error) {
	return FileMD5(path)
}

//计算文件的MD5值
func FileMD5(path string) (Md5 string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil)), err
}
