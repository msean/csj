package biz

import (
	"time"

	"gorm.io/gorm"
)

// customers表 结构体  Customers
type Customers struct {
	CreatedAt     time.Time      // 创建时间
	UpdatedAt     time.Time      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`                                                          // 删除时间
	Uid           string         `json:"uid" form:"uid" gorm:"primarykey;column:uid;comment:;size:64;"`           //uid字段
	OwnerUserUUID int64          `json:"-" form:"-" gorm:"column:owner_user;comment:'所属用户';size:64;"`             //'所属用户'
	Name          string         `json:"name" form:"name" gorm:"column:name;comment:'客户名字';size:4294967295;"`     //'客户名字'
	Phone         string         `json:"phone" form:"phone" gorm:"column:phone;comment:'手机号';size:4294967295;"`   //'手机号'
	Remark        string         `json:"remark" form:"remark" gorm:"column:remark;comment:'备注';size:4294967295;"` //'备注'
	Debt          *float64       `json:"debt" form:"debt" gorm:"column:debt;comment:'客户欠款,单位:元';size:10;"`        //'客户欠款,单位:元'
	Addr          string         `json:"addr" form:"addr" gorm:"column:addr;comment:'客户住址';size:4294967295;"`     //'客户住址'
	CarNo         string         `json:"carNo" form:"carNo" gorm:"column:car_no;comment:'车牌号';size:4294967295;"`  //'车牌号'
	OwnerUser     string         `json:"ownerUser" form:"-" gorm:"-"`                                             //'所属用户'
}

// TableName customers表 Customers自定义表名 customers
func (Customers) TableName() string {
	return "customers"
}
