package dao

import (
	"app/pkg/utils"
	"app/service/model"

	"gorm.io/gorm"
)

type batchGoodsDao struct{}

func newbatchGoodsDao() *batchGoodsDao {
	return &batchGoodsDao{}
}

func (dao *batchGoodsDao) FromUUID(db *gorm.DB, ownerUser, batchUUID, goodsUUID string) (*model.BatchGoods, error) {
	var batchGoods model.BatchGoods
	conds := []utils.Cond{
		utils.NewWhereCond("batch_uuid", batchUUID),
		utils.NewWhereCond("goods_uuid", goodsUUID),
		utils.NewWhereCond("owner_user", ownerUser),
	}

	if err := utils.Find(db, &batchGoods, conds...); err != nil {
		return nil, err
	}

	return &batchGoods, nil
}

func (dao *batchGoodsDao) ListByBatchUUID(db *gorm.DB, batchUUID string) (objects []model.BatchGoods, err error) {
	err = utils.Find(db, &objects, utils.NewWhereCond("batch_uuid", batchUUID))
	return
}
