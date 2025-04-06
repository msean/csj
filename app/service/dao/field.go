package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

func CustomerFieldSet(db *gorm.DB, uuid, ownerUser int64) (c model.CustomerField, err error) {
	var customer model.Customer
	if err = utils.GormFind(db, &customer, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	c.CustomerName = customer.Name
	c.CustomerPhone = customer.Phone
	return
}

func BatchCustomerFieldSet(db *gorm.DB, uuidList []int64, ownerUser int64) (customerM map[int64]model.CustomerField, err error) {
	customerM = make(map[int64]model.CustomerField)
	var customers []model.Customer
	if err = utils.GormFind(db, &customers, utils.NewWhereCond("owner_user", ownerUser), utils.NewInCondFromInt64("uid", uuidList)); err != nil {
		return
	}
	for _, c := range customers {
		customerM[c.UID] = model.CustomerField{
			CustomerName:  c.Name,
			CustomerPhone: c.Phone,
		}
	}
	return
}

func GoodsFieldSet(db *gorm.DB, uuid, ownerUser int64) (c model.GoodsField, err error) {
	var _goods model.Goods
	if err = utils.GormFind(db, &_goods, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	c = model.GoodsField{
		GoodsName: _goods.Name,
		GoodsTyp:  _goods.Typ,
	}
	return
}

func BatchGoodsFieldSet(db *gorm.DB, uuidList []int64, ownerUser int64) (goodsM map[int64]model.GoodsField, err error) {
	goodsM = make(map[int64]model.GoodsField)
	var _goodsList []model.Goods
	if err = utils.GormFind(db, &_goodsList, utils.NewWhereCond("owner_user", ownerUser), utils.NewInCondFromInt64("uid", uuidList)); err != nil {
		return
	}
	for _, c := range _goodsList {
		goodsM[c.UID] = model.GoodsField{
			GoodsName: c.Name,
			GoodsTyp:  c.Typ,
		}
	}
	return
}

// goodsM key goods uuid
func SetSurplusByBatch(db *gorm.DB, batchUUID int64) (goodsM map[int64]*model.SurplusField, err error) {
	goodsM = make(map[int64]*model.SurplusField)
	var _bgs []model.BatchGoods
	if err = utils.GormFind(db, &_bgs, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bg := range _bgs {
		goodsM[_bg.GoodsUUID] = &model.SurplusField{
			Weight: _bg.Weight,
			Mount:  _bg.Mount,
		}
	}
	var _bos []model.BatchOrderGoods
	if err = utils.GormFind(db, &_bos, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bo := range _bos {
		if _, ok := goodsM[_bo.GoodsUUID]; ok {

			goodsM[_bo.GoodsUUID] = &model.SurplusField{
				Weight: goodsM[_bo.GoodsUUID].Weight - _bo.Weight,
				Mount:  goodsM[_bo.GoodsUUID].Mount - _bo.Mount,
			}
		}
	}
	for _, surplus := range goodsM {
		surplus.Set()
	}
	return
}
