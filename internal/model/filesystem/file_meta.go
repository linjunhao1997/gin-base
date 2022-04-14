package filesystem

import (
	"gorm.io/gorm"
	"time"
)

type FileMeta struct {
	Id        int            `gorm:"column:id;primary_key" json:"id"`
	Name      string         `gorm:"column:name" json:"name"`
	Size      int64          `gorm:"column:size" json:"size"`
	Type      string         `gorm:"column:type" json:"type"`
	Bucket    string         `gorm:"column:bucket" json:"bucket"`
	CreatedAt time.Time      `gorm:"column:created_at;json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
