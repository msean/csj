package biz

import (
	"time"
)

type BaseModel struct {
	UID       int64      `gorm:"primaryKey;autoIncrement:false" json:"uuid"`
	CreatedAt time.Time  `json:"createTime"`
	UpdatedAt time.Time  `json:"updateTime"`
	DeletedAt *time.Time `json:"-"`
}
