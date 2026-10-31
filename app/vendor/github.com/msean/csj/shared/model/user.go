package model

import (
	"time"
)

type User struct {
	BaseModel
	Name               string    `gorm:"column:name;size:100;not null;comment:用户名" json:"name"`
	Phone              string    `gorm:"column:phone;unique;comment:手机号码" json:"phone"`
	ContactPhones      string    `gorm:"column:contact_phones;size:255;comment:分享联系方式(多个用逗号分隔)" json:"contactPhones"`
	Province           string    `gorm:"column:province;size:50;comment:省份" json:"province"`
	MarketName         string    `gorm:"column:market_name;size:100;comment:市场名称" json:"marketName"`
	MarketAddress      string    `gorm:"column:market_address;size:255;comment:市场地址" json:"marketAddress"`
	StallName          string    `gorm:"column:stall_name;size:100;comment:档口名称" json:"stallName"`
	ExperienceDeadline time.Time `gorm:"column:vip_expire_time;comment:体验截止时间" json:"vipRemainDays"`
	VipDeadline        time.Time `gorm:"column:vip_deadline;comment:会员截止时间" json:"vipDeadline"`
}

func (u *User) Title() string {
	if u.Name == "" {
		return u.Phone
	}
	return u.Name
}
