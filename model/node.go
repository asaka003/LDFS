package model

import "time"

type SampleUploadList struct {
	DataShards   int      `json:"data_shards"`
	ParityShards int      `json:"parity_shards"`
	UrlList      []string `json:"url_list"`
}

type SampleUploadInfo struct {
	FileKey string   `json:"file_key"`
	Shards  []*Shard `json:"shards"`
}

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
	Shards        []*Shard  `json:"shards"`
	EncodingTime  time.Time `json:"encoding_time"`
}

type Shard struct {
	ShardID   int      `json:"shard_id"`
	Hash      string   `json:"hash"`
	NodeName  string   `json:"node_name"`
	NodeURLs  []string `json:"node_urls"`
	PartHashs []string `json:"part_hashs"` // EC纠删码模式下可用
}

type CompleteMultipartUploadOptions struct {
	Parts []Object
}
type Object struct {
	Hash       string
	PartNumber int
}

type DataNode struct {
	// IP           string `json:"ip"`
	// Port         string `json:"port"`
	URL               string `json:"url"`
	NodeName          string `json:"node_name"`
	NodeDiskSize      int64  `json:"node_disk_size"`
	NodeFileTotalSize int64  `json:"node_file_total_size"`
	NodeDiskUsedSize  int64  `json:"node_disk_used_size"`
}
