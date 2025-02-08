package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	// mysql
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BaseModel struct {
	UID       string     `gorm:"primary_key;size:64" json:"uuid"`
	CreatedAt time.Time  `json:"createTime"`
	UpdatedAt time.Time  `json:"updateTime"`
	DeletedAt *time.Time `json:"-"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.UID == "" {
		id, _ := uuid.NewV4()
		m.UID = id.String()
	}
	return nil
}
