package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type BatchOrderGoodsDao struct{}

func NewBatchOrderGoodsDao() *BatchOrderGoodsDao {
	return &BatchOrderGoodsDao{}
}

func (dao *BatchOrderGoodsDao) Update(db *gorm.DB, batchOrderGoods model.BatchOrderGoods) error {
	return utils.WhereUIDCond(batchOrderGoods.UID).Cond(db).Updates(&model.BatchOrderGoods{
		Price:  batchOrderGoods.Price,
		Mount:  batchOrderGoods.Mount,
		Weight: batchOrderGoods.Weight,
	}).Error
}
