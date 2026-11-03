package model

import "time"

// QuickPayAudit retains the pre-payment snapshot and allocation evidence.
type QuickPayAudit struct {
	ID             uint64 `gorm:"primaryKey"`
	CreatedAt      time.Time
	OwnerUser      string  `gorm:"size:64;index"`
	CustomerUUID   string  `gorm:"size:64;index"`
	BeforeCredit   float64 `gorm:"type:decimal(18,2)"`
	BeforeDebt     float64 `gorm:"type:decimal(18,2)"`
	BeforePaidDebt float64 `gorm:"type:decimal(18,2)"`
	Received       float64 `gorm:"type:decimal(18,2)"`
	Discount       float64 `gorm:"type:decimal(18,2)"`
	Remaining      float64 `gorm:"type:decimal(18,2)"`
	Snapshot       string  `gorm:"type:longtext"`
	Allocations    string  `gorm:"type:longtext"`
}
