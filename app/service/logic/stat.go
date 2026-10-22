package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// goodsM key goods uuid
func BatchGoodsStat(db *gorm.DB, ownerUser string, batchUUID string) (goodsM map[string]*model.BatchGoodsStatFeild, err error) {
	goodsM = make(map[string]*model.BatchGoodsStatFeild)
	var _bgs []model.BatchGoods
	if err = utils.Find(db, &_bgs, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}

	for _, _bg := range _bgs {
		goodsM[_bg.GoodsUUID] = &model.BatchGoodsStatFeild{
			SurplusWeight:  _bg.Weight,
			SurplusMount:   _bg.Mount,
			Type:           _bg.GoodType,
			SellAmount:     "0",
			InventoryPrice: _bg.Price,
		}
	}

	var _bos []model.BatchOrderGoods
	err = db.Table(fmt.Sprintf("%s as bog", dao.ShardOrderGoodsTable(ownerUser))).
		Joins(fmt.Sprintf("JOIN %s AS bo ON bo.uid = bog.batch_order_uuid", dao.ShardOrderTable(ownerUser))).
		Where("bog.batch_uuid = ?", batchUUID).
		Where("bo.status IN (?)", common.ValidOrder).
		Find(&_bos).Error
	for _, _bo := range _bos {
		if _, ok := goodsM[_bo.GoodsUUID]; ok {
			global.Global.Logger.Debug("SetSurplusByBatch _bo",
				zap.Any("_bo.GoodsUUID", _bo.GoodsUUID),
				zap.Any("goodsM[_bo.GoodsUUID].Weight", goodsM[_bo.GoodsUUID].SurplusWeight),
				zap.Any("goodsM[_bo.GoodsUUID].Mount", goodsM[_bo.GoodsUUID].SurplusMount),
				zap.Any("_bo.Weight", _bo.Weight),
				zap.Any("_bo.Mount", _bo.Mount),
				zap.Any("_bo.Price", _bo.Price),
			)
			if goodsM[_bo.GoodsUUID].Type == common.GoodsTypeFix {
				goodsM[_bo.GoodsUUID] = &model.BatchGoodsStatFeild{
					SurplusMount:   goodsM[_bo.GoodsUUID].SurplusMount - _bo.Mount,
					Type:           goodsM[_bo.GoodsUUID].Type,
					InventoryPrice: goodsM[_bo.GoodsUUID].InventoryPrice,
					SellWeight:     goodsM[_bo.GoodsUUID].SellWeight + _bo.Weight,
					SellMount:      goodsM[_bo.GoodsUUID].SellMount + _bo.Mount,
					SellTotal:      goodsM[_bo.GoodsUUID].SellTotal + float64(_bo.Mount)*_bo.Price,
				}
			} else {
				goodsM[_bo.GoodsUUID] = &model.BatchGoodsStatFeild{
					SurplusWeight:  goodsM[_bo.GoodsUUID].SurplusWeight - _bo.Weight,
					Type:           goodsM[_bo.GoodsUUID].Type,
					InventoryPrice: goodsM[_bo.GoodsUUID].InventoryPrice,
					SellWeight:     goodsM[_bo.GoodsUUID].SellWeight + _bo.Weight,
					SellMount:      goodsM[_bo.GoodsUUID].SurplusMount + _bo.Mount,
					SellTotal:      goodsM[_bo.GoodsUUID].SellTotal + _bo.Weight*_bo.Price,
				}
			}
		}
	}

	global.Global.Logger.Debug("SetSurplusByBatch", zap.Any("goodsM", goodsM))
	for goodsUUID, surplus := range goodsM {
		surplus.Set()
		global.Global.Logger.Debug("#########SetSurplusByBatch",
			zap.Any("goods_uuid", zap.Any("goodsUUID", goodsUUID)),
			zap.Any("surplus Weight", surplus.SurplusWeight),
			zap.Any("surplus Mount", surplus.SurplusMount),
			zap.Any("surplus Type", surplus.Type),
			zap.Any("surplus SellTotal", surplus.SellTotal),
			zap.Any("surplus SellTotalAmount", surplus.SellTotalAmount),
			zap.Any("surplus SellAvgPrice", surplus.SellAvgPrice),
		)
	}
	return
}
