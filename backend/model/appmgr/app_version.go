package appmgr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AppVersion App版本管理
// 对应 app 客户端项目（app/service/model/version.go）的 app_versions 表，主键为 uid
type AppVersion struct {
	UID          string         `json:"uid" gorm:"column:uid;type:varchar(64);primaryKey"`                            // 主键
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at"`                                          // 创建时间
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"column:updated_at"`                                          // 更新时间
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`                                            // 删除时间
	Version      string         `json:"version" form:"version" gorm:"column:version;size:100;not null;comment:版本号;" binding:"required"`       // 版本号
	DownloadAddr string         `json:"downloadAddr" form:"downloadAddr" gorm:"column:download_addr;size:100;not null;comment:下载地址;" binding:"required"` // 下载地址
	IsForce      bool           `json:"isForce" form:"isForce" gorm:"column:is_force;default:false;comment:是否强制更新;"`                        // 是否强制更新
	Description  string         `json:"description" form:"description" gorm:"column:description;size:500;comment:更新说明;"`                    // 更新说明
}

// TableName AppVersion 自定义表名
func (AppVersion) TableName() string {
	return "app_versions"
}

// BeforeCreate 创建前自动生成主键 uid
func (m *AppVersion) BeforeCreate(tx *gorm.DB) error {
	if m.UID == "" {
		m.UID = uuid.NewString()
	}
	return nil
}
