package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchOrderPayDao struct{}

func NewBatchOrderPayDao() *BatchOrderPayDao {
	return &BatchOrderPayDao{}
}

func (dao *BatchOrderPayDao) Update(db *gorm.DB, batchOrderPay model.BatchOrderPay) error {
	return utils.WhereUIDCond(batchOrderPay.UID).Cond(db).Updates(&model.BatchOrderPay{
		PayType: batchOrderPay.PayType,
		Amount:  batchOrderPay.Amount,
	}).Error
}
