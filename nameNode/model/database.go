package model

import (
	"LDFS/model"
	"time"
)

/* 数据库结构 */

type File struct {
	Id         int                `db:"id"`
	FileName   string             `db:"file_name"`
	FileSize   int                `db:"file_size"`
	FileExt    string             `db:"file_ext"`
	Hash       string             `db:"hash"`
	MeteJson   string             `db:"meta"`
	Mete       model.FileMetadata `db:"-"`
	CreateTime time.Time          `db:"create_time"`
	UpdateTime time.Time          `db:"update_time"`
}

// type FileMetadata struct {
// 	UUID         string    `json:"uuid"`
// 	FileKey      string    `json:"file_key"`
// 	FileHash     string    `json:"file_hash"`
// 	ContentType  string    `json:"content-type"`
// 	Etag         string    `json:"etag"`
// 	FileSize     int64     `json:"file_size"`
// 	DataShards   int       `json:"data_shards"`
// 	ParityShards int       `json:"parity_shards"`
// 	Shards       []Shard   `json:"shards"`
// 	EncodingTime time.Time `json:"encoding_time"`
// }

// type Shard struct {
// 	ShardID  int    `json:"shard_id"`
// 	NodeName string `json:"node_name"`
// 	NodeURL  string `json:"node_url"`
// 	Hash     string `json:"hash"`
// }
