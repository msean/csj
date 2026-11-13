package model

// OrderTicketSettlement 保存收银时的金额，后续还款不会改写销售小票。
type OrderTicketSettlement struct {
	OwnerUser      string  `gorm:"primaryKey;size:64"`
	OrderUUID      string  `gorm:"primaryKey;size:64"`
	PreviousCredit float64 `gorm:"type:decimal(12,2);not null"`
	OrderCredit    float64 `gorm:"type:decimal(12,2);not null"`
	Received       float64 `gorm:"type:decimal(12,2);not null"`
	PayType        int32
}
