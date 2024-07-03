// 自动生成模板Users
package user

import (
	"time"

	"gorm.io/gorm"
)

// users表 结构体  Users
type Users struct {
	CreatedAt     time.Time      // 创建时间
	UpdatedAt     time.Time      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                                                                   // 删除时间
	Uid           string         `json:"uid" form:"uid" gorm:"primarykey;column:uid;comment:;size:64;"`                    //uid字段
	Name          string         `json:"name" form:"name" gorm:"column:name;comment:用户名;size:100;"`                        //用户名
	Phone         string         `json:"phone" form:"phone" gorm:"column:phone;comment:手机号码;size:191;"`                    //手机号码
	VipExpireTime *time.Time     `json:"vipExpireTime" form:"vipExpireTime" gorm:"column:vip_expire_time;comment:体验截止时间;"` //体验截止时间
}

// TableName users表 Users自定义表名 users
func (Users) TableName() string {
	return "users"
}
