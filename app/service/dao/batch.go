package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchDao struct{}

func NewBatchDao() *BatchDao {
	return &BatchDao{}
}

func (dao *BatchDao) Update(db *gorm.DB, batch model.Batch) (err error) {
	toUpdates := make(map[string]any)
	if batch.StorageTime != 0 {
		toUpdates["storage_time"] = batch.StorageTime
	}
	if batch.Status != 0 {
		toUpdates["status"] = batch.Status
	}
	if len(toUpdates) > 0 {
		return utils.WhereUIDCond(batch.UID).Cond(db).Model(&model.Batch{}).Updates(toUpdates).Error
	}
	return
}
