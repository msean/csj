package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// goodsM key goods uuid
func BatchGoodsStat(db *gorm.DB, batchUUID string) (goodsM map[string]*model.BatchGoodsStatFeild, err error) {
	goodsM = make(map[string]*model.BatchGoodsStatFeild)
	var _bgs []model.BatchGoods
	if err = utils.Find(db, &_bgs, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}

	for _, _bg := range _bgs {
		goodsM[_bg.GoodsUUID] = &model.BatchGoodsStatFeild{
			Weight:     _bg.Weight,
			Mount:      _bg.Mount,
			Type:       _bg.GoodType,
			SellAmount: "0",
		}
	}
	var _bos []model.BatchOrderGoods
	err = db.
		Table("batch_order_goods AS bog").
		Joins("JOIN batch_orders AS bo ON bo.uid = bog.batch_order_uuid").
		Where("bog.batch_uuid = ?", batchUUID).
		Where("bo.status IN (?)", common.ValidOrder).
		Find(&_bos).Error
	for _, _bo := range _bos {
		if _, ok := goodsM[_bo.GoodsUUID]; ok {
			global.Global.Logger.Debug("SetSurplusByBatch _bo",
				zap.Any("_bo.GoodsUUID", _bo.GoodsUUID),
				zap.Any("goodsM[_bo.GoodsUUID].Weight", goodsM[_bo.GoodsUUID].Weight),
				zap.Any("goodsM[_bo.GoodsUUID].Mount", goodsM[_bo.GoodsUUID].Mount),
				zap.Any("_bo.Weight", _bo.Weight),
				zap.Any("_bo.Mount", _bo.Mount),
				zap.Any("_bo.Price", _bo.Price),
			)
			if goodsM[_bo.GoodsUUID].Type == common.GoodsTypeFix {
				goodsM[_bo.GoodsUUID] = &model.BatchGoodsStatFeild{
					Weight:     goodsM[_bo.GoodsUUID].Weight - _bo.Weight,
					Mount:      goodsM[_bo.GoodsUUID].Mount - _bo.Mount,
					Type:       goodsM[_bo.GoodsUUID].Type,
					SellWeight: goodsM[_bo.GoodsUUID].SellWeight + _bo.Weight,
					SellMount:  goodsM[_bo.GoodsUUID].SellMount + _bo.Mount,
					SellTotal:  goodsM[_bo.GoodsUUID].SellTotal + float64(_bo.Mount)*_bo.Price,
				}
			} else {
				goodsM[_bo.GoodsUUID] = &model.BatchGoodsStatFeild{
					Weight:     goodsM[_bo.GoodsUUID].Weight - _bo.Weight,
					Type:       goodsM[_bo.GoodsUUID].Type,
					SellWeight: goodsM[_bo.GoodsUUID].SellWeight + _bo.Weight,
					SellMount:  goodsM[_bo.GoodsUUID].Mount + _bo.Mount,
					SellTotal:  goodsM[_bo.GoodsUUID].SellTotal + _bo.Weight*_bo.Price,
				}
			}
		}
	}

	global.Global.Logger.Info("SetSurplusByBatch", zap.Any("goodsM", goodsM))
	for goodsUUID, surplus := range goodsM {
		surplus.Set()
		global.Global.Logger.Info("SetSurplusByBatch",
			zap.Any("goods_uuid", zap.Any("goodsUUID", goodsUUID)),
			zap.Any("surplus Weight", surplus.Weight),
			zap.Any("surplus Mount", surplus.Mount),
			zap.Any("surplus Surplus", surplus.Surplus),
			zap.Any("surplus Type", surplus.Type),
			zap.Any("surplus SellTotal", surplus.SellTotal),
			zap.Any("surplus SellTotalAmount", surplus.SellTotalAmount),
			zap.Any("surplus SellAvgPrice", surplus.SellAvgPrice),
		)
	}
	return
}
