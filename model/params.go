package model

type DownloadShardParam struct {
	Hash string `json:"hash"`
}

type RecoverShardParam struct {
	Block          *Block `json:"block"`
	ShardId        int64  `json:"shard_id"`
	DataShardNum   int    `json:"data_shard_num"`   //EC模式下可用
	ParityShardNum int    `json:"parity_shard_num"` //EC模式下可用
	StoragePolicy  string `json:"storage_policy"`
}

type RequestUploadFileParams struct {
	FileKey       string `json:"file_key"`
	FileSize      int64  `json:"file_size"`
	BlockSize     int64  `json:"block_size"`
	StoragePolicy string `json:"storage_policy"`
}

type CompleteSampleUploadParams struct {
	FileKey string `json:"file_key"`
}
