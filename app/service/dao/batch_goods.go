package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchGoodsDao struct{}

func NewBatchGoods() *BatchGoodsDao {
	return &BatchGoodsDao{}
}

func (dao *BatchGoodsDao) Update(db *gorm.DB, batchGoods model.BatchGoods) error {
	return utils.WhereUIDCond(batchGoods.UID).Cond(db).Updates(&model.BatchGoods{
		Price:  batchGoods.Price,
		Weight: batchGoods.Weight,
		Mount:  batchGoods.Mount,
	}).Error
}
