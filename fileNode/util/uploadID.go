package util

/*
	生成用户分块上传ID值
*/

import (
	"LDFS/fileNode/web_pkg/snowflake"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var (
	ErrUploadID = errors.New("UploadID must be a 64-bit hexadecimal string")
)

func GenerateUploadID(file_key string) string {
	id := snowflake.GetID()
	// 将 ID 转换为字节序列
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(id >> ((7 - i) * 8))
	}

	// 计算 SHA-256 哈希值
	h := sha256.Sum256(append(b, []byte(file_key)...))

	// 将哈希值转换为 64 位的十六进制字符串
	hashedID := hex.EncodeToString(h[:])
	return hashedID
}

//验证是否为UploadID
func IsUploadID(UploadID string) bool {
	if len(UploadID) != 64 {
		return false
	}
	_, err := hex.DecodeString(UploadID)
	return err == nil
}
