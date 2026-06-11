package dao

import (
	"app/pkg/utils"
	"app/service/model"

	"gorm.io/gorm"
)

type batchDao struct{}

func newbatchDao() *batchDao {
	return &batchDao{}
}

func (dao *batchDao) Update(db *gorm.DB, object model.Batch) (err error) {
	toUpdates := make(map[string]any)
	if object.StorageTime != 0 {
		toUpdates["storage_time"] = object.StorageTime
	}
	if object.Status != 0 {
		toUpdates["status"] = object.Status
	}
	if len(toUpdates) > 0 {
		return utils.WhereUIDCond(object.UID).Cond(db).Model(&model.Batch{}).Updates(toUpdates).Error
	}
	return
}

func (dao *batchDao) UpdateBatchGoods(db *gorm.DB, object model.BatchGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchGoods{
		Price:  object.Price,
		Weight: object.Weight,
		Mount:  object.Mount,
	}).Error
}
