package model

import (
	"time"
)

type File struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	FileSize   int64
	FileKey    string `gorm:"index"`
	FileExt    string
	Hash       string `gorm:"index"`
	Meta       string
	CreateTime *time.Time
	UpdateTime *time.Time
}
