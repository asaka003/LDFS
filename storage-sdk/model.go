package storagesdk

type InitiateMultipartUploadRequest struct {
	FileKey string `json:"file_key"`
}

type InitiateMultipartUploadResponse struct {
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type ListPartsRequest struct {
	UploadID string `json:"upload_id"`
}
