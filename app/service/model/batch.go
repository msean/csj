package model

import (
	"app/service/common"

	"gorm.io/gorm"
)

type (
	Batch struct {
		BaseModel
		OwnerUser        string        `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		SerialNo         string        `gorm:"column:serial_no;comment:序号" json:"serialNo"`
		StorageTime      int64         `gorm:"column:storage_time;comment:入库时间;type:bigint" json:"storageTime"`
		Status           int32         `gorm:"column:status;comment:状态" json:"status"`
		GoodsListRelated []*BatchGoods `gorm:"foreignKey:BatchUUID;references:UID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"goodsList"`
		SerialID         int           `gorm:"column:serial_id;comment:ownerUserr自增ID" json:"serialID"`
		BatchStatFeild
	}
	BatchGoods struct {
		BaseModel
		OwnerUser string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		BatchUUID string  `gorm:"column:batch_uuid;comment:货品uuid" json:"batchUUID"`
		SerialNo  string  `gorm:"column:serial_no;comment:批次序号" json:"serialNo"`
		GoodsUUID string  `gorm:"column:goods_uuid;comment:货品uuid" json:"goodsUUID"`
		GoodType  int     `gorm:"column:goods_type;comment:货品类型" json:"goods_type"`
		Price     float64 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
		Weight    float64 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
		Mount     int     `gorm:"column:mount;comment:数量" json:"mount"`
		GoodsFeild
		BatchGoodsStatFeild
	}
)

// BeforeCreate 中按用户自动递增 ID
func (b *Batch) BeforeCreate(tx *gorm.DB) (err error) {
	b.BaseModel.BeforeCreate(tx)
	b.Default()
	if b.SerialID == 0 {
		if b.SerialID, err = b.GenerateSerialID(tx); err != nil {
			return
		}
	}
	// if b.SerialNo == "" {
	// 	b.SerialNo = fmt.Sprintf("%s-%d", b.OwnerUser, b.SerialID)
	// }

	return nil
}

func (Batch) TableName() string { return "batches" }

func (b *Batch) GenerateSerialID(tx *gorm.DB) (maxID int, err error) {
	if err = tx.Model(&Batch{}).
		Where("owner_user = ?", b.OwnerUser).
		Select("COALESCE(MAX(serial_id), 0)").
		Scan(&maxID).Error; err != nil {
		return
	}
	maxID = maxID + 1
	return
}

func (b *Batch) Default() {
	b.Status = common.BatchStatusOnSellering
	// b.SerialNo = SerioalNo(time.Time{})
}

// func (bg *BatchGoods) Amount() (amount string) {
// 	if bg.Type == common.GoodsTypeFix {
// 		amount = fmt.Sprintf("%d", bg.Mount)
// 		return
// 	}
// 	amount = fmt.Sprintf("%.2f", bg.Weight)
// 	return
// }
