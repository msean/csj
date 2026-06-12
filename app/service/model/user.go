package model

import (
	"app/pkg/utils"
	"time"
)

type User struct {
	BaseModel
	Name               string     `gorm:"column:name;size:100;not null;comment:用户名" json:"name"`
	Phone              string     `gorm:"column:phone;unique;comment:手机号码" json:"phone"`
	ExperienceDeadline *time.Time `gorm:"column:vip_expire_time;comment:体验截止时间" json:"vipRemainDays"`
}

func (u *User) Title() string {
	if utils.IsBlankString(u.Name) {
		return u.Phone
	}
	return u.Name
}
