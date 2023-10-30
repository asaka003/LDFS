package util

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
)

type Sha1Stream struct {
	_sha1 hash.Hash
}

func (obj *Sha1Stream) Update(data []byte) {
	if obj._sha1 == nil {
		obj._sha1 = sha1.New()
	}
	obj._sha1.Write(data)
}

func (obj *Sha1Stream) Sum() string {
	return hex.EncodeToString(obj._sha1.Sum([]byte("")))
}

func GetMutiFileSha1(file *multipart.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, *file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

//支持大文件计算hash值，只占用32K左右内存
/*
	hash接口实现了一个write的方法，
	会计算每次写入数据后的hash值，
	IO.COPY每次循环调用写入数据的时候都会调用这个write方法，
	当数据读完以后hash值也就计算出来了，
	本质上还是从文件句柄里面一次次读取数据然后计算hash值，
	所以占用内存就只有32K以上，不会占用大内存。
*/

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

func IsMD5(s string) bool {
	// 匹配MD5值的正则表达式
	r := regexp.MustCompile(`^[a-f0-9]{32}$`)

	// 返回字符串是否符合正则表达式
	return r.MatchString(s)
}

func IsUUID(s string) bool {
	// UUID正则表达式，匹配8-4-4-12格式的UUID
	pattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	uuidRegex := regexp.MustCompile(pattern)

	return uuidRegex.MatchString(s)
}

func CreateDir(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func RecreateDir(dir string) error {
	_, err := os.Stat(dir)
	if err == nil { //目录存在
		if err = os.RemoveAll(dir); err != nil {
			return err
		}
		CreateDir(dir)
	} else if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return err
}
