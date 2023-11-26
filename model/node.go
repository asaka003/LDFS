package model

import "time"

type FileMetadata struct {
	UUID          string    `json:"uuid"`
	FileKey       string    `json:"file_key"`
	FileHash      string    `json:"file_hash"`
	ContentType   string    `json:"content-type"`
	Etag          string    `json:"etag"`
	FileSize      int64     `json:"file_size"`
	DataShards    int       `json:"data_shards"`
	ParityShards  int       `json:"parity_shards"`
	StoragePolicy string    `json:"storage_policy"` //存储策略  EC表示纠删码模式,copy表示副本冗余模式
	Blocks        []*Block  `json:"blocks"`
	CreateTime    time.Time `json:"create_time"`
	Status        string    `json:"status"`
}

type Block struct {
	BlockId   int      `json:"block_id"`
	BlockSize int64    `json:"block_size"`
	Hash      string   `json:"hash"`
	Shards    []*Shard `json:"shards"`
}

type Shard struct {
	ShardID  int64  `json:"shard_id"`
	Hash     string `json:"hash"`
	NodeName string `json:"node_name"`
	NodeURL  string `json:"node_url"`
}

type DataNode struct {
	// IP           string `json:"ip"`
	// Port         string `json:"port"`
	URL                   string `json:"url"`
	NodeName              string `json:"node_name"`
	NodeDiskSize          int64  `json:"node_disk_size"`
	NodeFileTotalSize     int64  `json:"node_file_total_size"`
	NodeDiskUsedSize      int64  `json:"node_disk_used_size"`
	NodeDiskAvailableSize int64  `json:"node_disk_available_size"`
}

type FileInfo struct {
	FileKey string `json:"file_key"`
	Size    int64  `json:"size"`
}
