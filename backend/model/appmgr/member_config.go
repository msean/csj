package appmgr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberConfig 会员收费配置
// 对应 app 客户端项目（app/service/model/member_config.go）的 member_configs 表，主键为 uid
// type 唯一，取值 1月 2季度 3年
type MemberConfig struct {
	UID       string         `json:"uid" gorm:"column:uid;type:varchar(64);primaryKey"`                                                   // 主键
	CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`                                                                  // 创建时间
	UpdatedAt time.Time      `json:"updatedAt" gorm:"column:updated_at"`                                                                  // 更新时间
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`                                                                    // 删除时间
	Typ       int            `json:"type" form:"type" gorm:"column:type;uniqueIndex;not null;comment:会员类型 1月 2季度 3年;" binding:"required"` // 会员类型
	Price     float64        `json:"price" form:"price" gorm:"column:price;type:decimal(10,2);comment:收费金额;" binding:"required"`          // 收费金额
}

// TableName MemberConfig 自定义表名
func (MemberConfig) TableName() string {
	return "member_configs"
}

// BeforeCreate 创建前自动生成主键 uid
func (m *MemberConfig) BeforeCreate(tx *gorm.DB) error {
	if m.UID == "" {
		m.UID = uuid.NewString()
	}
	return nil
}
