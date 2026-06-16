package model

type BatchSerialCounter struct {
	OwnerUser   string `gorm:"column:owner_user;primaryKey;size:64" json:"ownerUser"`
	MaxSerialID int    `gorm:"column:max_serial_id;default:0" json:"maxSerialID"`
}

func (BatchSerialCounter) TableName() string { return "batch_serial_counters" }
