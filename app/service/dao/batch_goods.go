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
	return utils.WhereUIDCond(batchGoods.UID).Cond(db.Table(batchGoods.TableName())).Updates(map[string]any{
		"price":  batchGoods.Price,
		"weight": batchGoods.Weight,
		"mount":  batchGoods.Mount,
	}).Error
}

func (logic *BatchGoodsDao) FromBatchUUID(db *gorm.DB, batchUUID int64) (batchGoods []model.BatchGoods, err error) {
	err = utils.GormFind(db, &batchGoods, utils.NewWhereCond("batch_uuid", batchUUID))
	return
}
