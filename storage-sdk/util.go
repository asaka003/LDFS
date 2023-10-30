package storagesdk

import (
	"math/rand"
	"time"
)

// 生成指定长度的随机字符串
func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	// 定义包含可用字符的字符集
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	// 用于存储随机生成的字符
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		// 从字符集中随机选择一个字符
		randomIndex := rand.Intn(len(charset))
		result[i] = charset[randomIndex]
	}

	return string(result)
}

//计算文件的hash值
