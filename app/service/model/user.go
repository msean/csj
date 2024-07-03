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
	// Customers          []Customer `gorm:"foreignkey:OwnerUser; references:UID"`
}

func (u *User) Update(db *gorm.DB) (err error) {
	toUpdate := make(map[string]any)
	if u.Name != "" {
		toUpdate["name"] = u.Name
	}
	if u.Phone != "" {
		toUpdate["phone"] = u.Phone
	}
	return WhereUIDCond(u.UID).Cond(db.Model(&User{})).Updates(toUpdate).Error
}

func (u *User) WherePhoneCond() (w WhereCond) {
	return NewWhereCond("phone", u.Phone)
}
