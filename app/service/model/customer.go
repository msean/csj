package model

type Customer struct {
	BaseModel
	OwnerUser string  `gorm:"column:owner_user;comment:所属用户;type:VARCHAR(64)" json:"ownerUser"`
	Name      string  `gorm:"column:name;comment:客户名字" json:"name"`
	Phone     string  `gorm:"column:phone;comment:手机号" json:"phone"`
	Remark    string  `gorm:"column:remark;comment:备注" json:"remark"`
	Debt      float64 `gorm:"column:debt;type:decimal(10,2);comment:客户欠款,单位:元" json:"debt"`
	Addr      string  `gorm:"column:addr;comment:客户住址" json:"addr"`
	CarNo     string  `gorm:"column:car_no;comment:车牌号" json:"carNo"`
}
