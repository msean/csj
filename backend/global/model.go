package global

import (
	"time"

	"gorm.io/gorm"
)

type GVA_MODEL struct {
	ID        uint           `gorm:"primarykey" json:"ID"`                                 // 主键ID
	CreatedAt time.Time      `json:"createdAt" form:"createdAt" gorm:"column:created_at;"` // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                       // 删除时间
}
