package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name               string     `gorm:"column:name;size:100;not null;comment:用户名" json:"name"`
	Phone              string     `gorm:"column:phone;unique;comment:手机号码" json:"phone"`
	ExperienceDeadline *time.Time `gorm:"column:vip_expire_time;comment:体验截止时间" json:"vipRemainDays"`
}

func (User) TableName() string {
	return "users"
}

// 在插入前自动生成 ID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UID == 0 {
		// 为 "users" 表创建单独的 Snowflake 生成器（假设 nodeID=1）
		generator := utils.GetOrCreateGenerator("users", 1)
		u.UID = generator.GenerateID()
	}
	return nil
}
