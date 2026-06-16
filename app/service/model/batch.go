package model

import (
	"app/service/common"
	"strings"

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
	// MySQL atomic increment using LAST_INSERT_ID
	result := tx.Exec(`
		UPDATE batch_serial_counters
		SET max_serial_id = LAST_INSERT_ID(max_serial_id + 1)
		WHERE owner_user = ?
	`, b.OwnerUser)
	if result.Error != nil {
		err = result.Error
		return
	}

	if result.RowsAffected > 0 {
		// Existing counter updated, get the new value
		err = tx.Raw("SELECT LAST_INSERT_ID()").Scan(&maxID).Error
		return
	}

	// No counter exists yet — insert with initial value 1
	if err = tx.Exec(`
		INSERT INTO batch_serial_counters (owner_user, max_serial_id)
		VALUES (?, 1)
	`, b.OwnerUser).Error; err != nil {
		// Ignore duplicate key (another request may have just inserted)
		if !isDuplicateKeyError(err) {
			return
		}
		// Another request won the race, retry the update
		result = tx.Exec(`
			UPDATE batch_serial_counters
			SET max_serial_id = LAST_INSERT_ID(max_serial_id + 1)
			WHERE owner_user = ?
		`, b.OwnerUser)
		if result.Error != nil {
			err = result.Error
			return
		}
		err = tx.Raw("SELECT LAST_INSERT_ID()").Scan(&maxID).Error
		return
	}

	maxID = 1
	return
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// Check for MySQL duplicate key error (Error 1062)
	return strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "1062")
}

func (b *Batch) Default() {
	b.Status = common.BatchStatusOnSellering
	// b.SerialNo = SerioalNo(time.Time{})
}
