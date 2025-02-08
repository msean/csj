package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchOrderDao struct{}

func NewBatchOrderDao() *BatchOrderDao {
	return &BatchOrderDao{}
}

func (dao *BatchOrderDao) UpdateStatus(db *gorm.DB, uid string, status int32) error {
	return utils.WhereUIDCond(uid).Cond(db).Model(&model.BatchOrder{}).Update("status", status).Error
}

func (dao *BatchOrderDao) UpdateShare(db *gorm.DB, uid string, share int32) error {
	return utils.WhereUIDCond(uid).Cond(db).Updates(&model.BatchOrder{
		Shared: model.BatchOrderShared,
	}).Error
}
