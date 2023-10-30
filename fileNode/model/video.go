package model

type VideoInfo struct {
	FileID     int64  `db:"file_hash"`
	FaceImgUrl string `db:"face_url"`
}
