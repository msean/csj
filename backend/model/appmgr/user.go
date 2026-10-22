package appmgr

import (
	"time"

	"gorm.io/gorm"
)

// AppUser 对应 app 客户端 users 表的注册用户
// 表名由 app/service/model 的 User 模型决定（默认 users）
type AppUser struct {
	UID                string         `json:"uid" gorm:"column:uid;type:varchar(64);primaryKey"`
	CreatedAt          time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
	Name               string         `json:"name" gorm:"column:name;size:100;comment:用户名"`
	Phone              string         `json:"phone" gorm:"column:phone;type:varchar(32);uniqueIndex;comment:手机号码"`
	ExperienceDeadline time.Time      `json:"experienceDeadline" gorm:"column:vip_expire_time;comment:体验截止时间"`
	VipDeadline        time.Time      `json:"vipDeadline" gorm:"column:vip_deadline;comment:会员截止时间"`
}

func (AppUser) TableName() string {
	return "users"
}
