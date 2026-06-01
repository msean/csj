package model

import (
	"app/service/common"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	BatchStatusOnSellering = 1 // 在售
	BatchStatusOffShelf    = 2 // 售罄
	BatchStatusSettled     = 3 // 已结算

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
		SurplusFeild
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

func (b *Batch) Update(db *gorm.DB) (err error) {
	toUpdates := make(map[string]any)
	if b.StorageTime != 0 {
		toUpdates["storage_time"] = b.StorageTime
	}
	if b.Status != 0 {
		toUpdates["status"] = b.Status
	}
	if len(toUpdates) > 0 {
		return WhereUIDCond(b.UID).Cond(db).Model(&Batch{}).Updates(toUpdates).Error
	}
	return
}

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

func SerioalNo(dt time.Time) string {
	if dt.IsZero() {
		dt = time.Now()
	}
	return fmt.Sprintf("%d-%d-%d", dt.Year(), dt.Month(), dt.Day())
}

func (b *Batch) Default() {
	b.Status = BatchStatusOnSellering
	b.SerialNo = SerioalNo(time.Time{})
}

func (b *BatchOrder) DefaultSet() {
	b.Status = common.BatchOrderTemp
	b.Shared = common.BatchOrderUnshare
}

func (bg *BatchGoods) Update(db *gorm.DB) error {
	return WhereUIDCond(bg.UID).Cond(db).Updates(&BatchGoods{
		Price:  bg.Price,
		Weight: bg.Weight,
		Mount:  bg.Mount,
	}).Error
}

func (bg *BatchGoods) Amount() (amount string) {
	if bg.Type == common.GoodsTypeFix {
		amount = fmt.Sprintf("%d", bg.Mount)
		return
	}
	amount = fmt.Sprintf("%.2f", bg.Mount)
	return
}
