package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
)

//获取文件hash值
func MutiPartFileHash(file *multipart.FileHeader) (hash string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()
	md5_hash := md5.New()
	// 分块读取文件并逐次计算MD5值,优化内存占用
	buf := make([]byte, 1024*1024*10) // 10MB 缓冲区
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		md5_hash.Write(buf[:n]) //每次write都会计算一次MD5
	}
	md5sum := md5_hash.Sum(nil)
	hash = hex.EncodeToString(md5sum)
	return
}

func FileHash(file *os.File) string {
	return fileMD5(file)
}

func BytesHash(data []byte) string {
	return bytesMD5(data)
}

func bytesSha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func fileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func bytesMD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum(nil))
}

func fileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}
