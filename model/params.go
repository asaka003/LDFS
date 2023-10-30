package model

type DownloadShardParam struct {
	Hash string `json:"hash"`
}

type InitUploadParam struct {
	FileKey  string `json:"file_key"`
	FileHash string `json:"file_hash"`
	FileSize int64  `json:"file_size"`
}

type CompleteMultipartParam struct {
	FileKey  string                          `json:"file_key"`
	UploadID string                          `json:"upload_id"`
	FileHash string                          `json:"file_hash"`
	Opt      *CompleteMultipartUploadOptions `json:"opt"`
}

type AbortMultipartUploadParam struct {
	FileKey  string `json:"file_key"`
	UploadID string `json:"upload_id"`
}

type ListPartsParam struct {
	FileKey  string `json:"file_key"`
	UploadID string `json:"upload_id"`
}

type RequestUploadFileParams struct {
	FileKey       string `json:"file_key"`
	FileSize      int64  `json:"file_size"`
	BlockSize     int64  `json:"block_size"`
	StoragePolicy string `json:"storage_policy"`
}
