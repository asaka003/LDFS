package util

import (
	"LDFS/model"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

//生成指定长度的数字验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func IsMD5(s string) bool {
	// 匹配MD5值的正则表达式
	r := regexp.MustCompile(`^[a-f0-9]{32}$`)

	// 返回字符串是否符合正则表达式
	return r.MatchString(s)
}

func BytesHash(data []byte) string {
	return bytesMD5(data)
}

func bytesMD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum(nil))
}

//读取文件中的fileMeta信息
func GetFileMetaInFile(path string) (meta *model.FileMetadata, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	metabuf := make([]byte, 1024)
	buffer := bytes.NewBuffer(metabuf)
	io.Copy(buffer, file)
	meta = new(model.FileMetadata)
	err = json.Unmarshal(buffer.Bytes(), meta)
	return
}
