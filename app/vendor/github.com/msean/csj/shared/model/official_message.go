package model

import "time"

// OfficialMessage is shared by all users; publishing never fans out into MessageCenter.
type OfficialMessage struct {
	BaseModel
	Level       int        `json:"level" gorm:"default:1"` // 1 general, 2 important
	Title       string     `json:"title" gorm:"size:200"`
	Content     string     `json:"content" gorm:"type:text"`
	Status      int        `json:"status" gorm:"index;default:0"` // 0 draft, 1 published
	PublishedAt time.Time  `json:"publishedAt" gorm:"index"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	Sort        int        `json:"sort" gorm:"default:0"`
}

type OfficialMessageRead struct {
	UserUUID    string `gorm:"size:64;primaryKey"`
	MessageUUID string `gorm:"size:64;primaryKey"`
	ReadAt      time.Time
}
