package model

type ListPartsResponse struct {
	Msg   string    `json:"msg"`
	Parts []*Object `json:"parts"`
}

type RequestUploadFileResponse struct {
	FileMeta *FileMetadata `json:"file_meta"`
}
