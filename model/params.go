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

// type ECconfig struct {
// 	DataShardNum    int
// 	ParityShardsNum int
// }

type RequestUploadFileParams struct {
	FileKey       string `json:"file_key"`
	FileSize      int64  `json:"file_size"`
	BlockSize     int64  `json:"block_size"` //设置的每个块的最大值
	StoragePolicy string `json:"storage_policy"`
	// ECconfig      *ECconfig `json:"EC_config"` //纠删码模式下有效
}

type CompleteSampleUploadParams struct {
	FileKey  string        `json:"file_key"`
	FileMeta *FileMetadata `json:"file_meta"`
}

const (
	TypeDataNode string = "data-node"
	TypeNameNode string = "Name-node"
)

type ParamJoin struct {
	RaftAddr string `json:"addr"`
	HttpAddr string `json:"haddr"`
	ID       string `json:"id"`
}

type ParamJoinDataNode struct {
	DataNodeInfo *DataNode `json:"data_node"`
}

type ParamReplicasDownloadShard struct {
	FileKey string `json:"file_key"`
	BlockId int    `json:"block_id"`
	Shard   *Shard `json:"shard"`
}
