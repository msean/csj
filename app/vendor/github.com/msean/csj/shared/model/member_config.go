package model

// MemberConfig 会员收费配置
// type 唯一，取值见 common 中的 MemberTypeMonth / MemberTypeQuarter / MemberTypeYear
type MemberConfig struct {
	BaseModel
	Typ   int     `gorm:"column:type;comment:会员类型 1月 2季度 3年;uniqueIndex" json:"type"`
	Price float64 `gorm:"column:price;type:decimal(10,2);comment:收费金额" json:"price"`
}

func (MemberConfig) TableName() string {
	return "member_configs"
}
