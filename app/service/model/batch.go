package model

import (
	"fmt"
	"time"
)

const (
	BatchStatusOnSelling = 1 // 在售
	BatchStatusOffShelf  = 2 // 售罄
	BatchStatusSettled   = 3 // 已结算

)

type (
	Batch struct {
		BaseModel
		OwnerUser   int64  `gorm:"column:owner_user;comment:所属用户" json:"-"`
		SerialNo    string `gorm:"column:serial_no;comment:序号" json:"serialNo"`
		StorageTime int64  `gorm:"column:storage_time;comment:入库时间;type:bigint" json:"storageTime"`
		Status      int32  `gorm:"column:status;comment:状态" json:"status"`
	}
	BatchGoods struct {
		BaseModel
		OwnerUser int64   `gorm:"column:owner_user;comment:所属用户" json:"-"`
		BatchUUID int64   `gorm:"column:batch_uuid;comment:货品uuid" json:"-"`
		SerialNo  string  `gorm:"column:serial_no;comment:批次序号" json:"serialNo"`
		GoodsUUID int64   `gorm:"column:goods_uuid;comment:货品uuid" json:"-"`
		Price     float64 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
		Weight    float64 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
		Mount     int32   `gorm:"column:mount;comment:数量" json:"mount"`
	}
)

func SerialNo(dt time.Time) string {
	if dt.IsZero() {
		dt = time.Now()
	}
	return fmt.Sprintf("%d-%02d-%02d", dt.Year(), dt.Month(), dt.Day())
}

func (b *Batch) Default() {
	b.Status = BatchStatusOnSelling
	b.SerialNo = SerialNo(time.Time{})
}

func (b *BatchOrder) DefaultSet() {
	b.Status = BatchOrderTemp
	b.Shared = BatchOrderUnshare
}

func (BatchGoods) TableName() string {
	return "batch_goods"
}
