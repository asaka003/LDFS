package util

import (
	"LDFS/dataNode/util"
	"LDFS/model"
	"LDFS/nameNode/raft"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
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
func GetFileMetaInFile(key string) (meta *model.FileMetadata, err error) {
	if key[0] == '/' {
		key = key[1:]
	}
	key_path := util.BytesHash([]byte(key)) + ".json"
	rmeta, err := raft.RaftNodeClient.GetFileMeta(key_path)
	meta = (*model.FileMetadata)(rmeta)
	return
}

//存储文件meta信息到文件中
func SaveFileMetaInFile(fileMeta *model.FileMetadata) (err error) {
	//保存meta信息到文件中
	return raft.RaftNodeClient.CreateFileMeta(util.BytesHash([]byte(fileMeta.FileKey))+".json", (*raft.FileMeta)(fileMeta))
}

//更新文件meta信息到文件中
func UpdateFileMeta(fileMeta *model.FileMetadata) (err error) {
	return raft.RaftNodeClient.UpdateFileMeta(util.BytesHash([]byte(fileMeta.FileKey))+".json", (*raft.FileMeta)(fileMeta))
}

//删除文件meta信息
func DeleteFileMeta(key string) (err error) {
	key_path := util.BytesHash([]byte(key)) + ".json"
	err = raft.RaftNodeClient.DeleteFileMeta(key_path)
	return
}

//获取leader节点的URL
func GetNameNodeLeaderInfo() (leader *model.NameNode, err error) {
	leader = raft.RaftNodeClient.GetLeaderNameNode()
	if leader.HAddr == "" {
		err = errors.New("正在选举leader,请稍后重试")
	}
	return
}
