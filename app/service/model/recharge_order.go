package model

import "time"

// RechargeOrder 充值订单（会员充值）
type RechargeOrder struct {
	BaseModel
	UserUUID      string     `gorm:"column:user_uuid;type:varchar(64);index;comment:用户UUID" json:"userUUID"`
	OrderNo       string     `gorm:"column:order_no;type:varchar(64);uniqueIndex;comment:商户订单号" json:"orderNo"`
	MemberType    int        `gorm:"column:member_type;comment:会员类型 1月 2季度 3年" json:"memberType"`
	Amount        float64    `gorm:"column:amount;type:decimal(10,2);comment:充值金额" json:"amount"`
	PayStatus     int        `gorm:"column:pay_status;default:0;index;comment:支付状态 0待支付 1已支付 2已取消 3已退款" json:"payStatus"`
	PayTime       *time.Time `gorm:"column:pay_time;comment:支付时间" json:"payTime"`
	TransactionID string     `gorm:"column:transaction_id;type:varchar(64);comment:微信支付订单号" json:"transactionID"`
}

func (RechargeOrder) TableName() string {
	return "recharge_orders"
}
