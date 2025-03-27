package model

import (
	"time"
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
