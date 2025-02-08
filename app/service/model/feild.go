package model

import (
	"app/utils"
	"fmt"

	"gorm.io/gorm"
)

type (
	SurplusField struct {
		Weight  float64 `gorm:"-" json:"-"`
		Mount   int32   `gorm:"-" json:"-"`
		Surplus string  `gorm:"-" json:"surplus"`
	}
	CustomerField struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsField struct {
		GoodsName string `gorm:"-" json:"name"`
		GoodsTyp  int32  `gorm:"-" json:"type"`
	}
	PayField struct {
		PayFee  float64 `gorm:"-"  json:"payFee"`
		PayType int32   `gorm:"-"  json:"payType"`
		PaidFee float64 `gorm:"-"  json:"paidFee"`
	}
)

func CustomerFieldSet(db *gorm.DB, uuid, ownerUser string) (c CustomerField, err error) {
	var customer Customer
	if err = utils.GormFind(db, &customer, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	c.CustomerName = customer.Name
	c.CustomerPhone = customer.Phone
	return
}

func BatchCustomerFieldSet(db *gorm.DB, uuidList []string, ownerUser string) (customerM map[string]CustomerField, err error) {
	customerM = make(map[string]CustomerField)
	var customers []Customer
	if err = utils.GormFind(db, &customers, utils.NewWhereCond("owner_user", ownerUser), utils.NewInCondFromString("uid", uuidList)); err != nil {
		return
	}
	for _, c := range customers {
		customerM[c.UID] = CustomerField{
			CustomerName:  c.Name,
			CustomerPhone: c.Phone,
		}
	}
	return
}

func GoodsFeildSet(db *gorm.DB, uuid, ownerUser string) (c GoodsField, err error) {
	var _goods Goods
	if err = utils.GormFind(db, &_goods, utils.NewWhereCond("owner_user", ownerUser), utils.WhereUIDCond(uuid)); err != nil {
		return
	}
	c = GoodsField{
		GoodsName: _goods.Name,
		GoodsTyp:  _goods.Typ,
	}
	return
}

func BatchGoodsFeildSet(db *gorm.DB, uuidList []string, ownerUser string) (goodsM map[string]GoodsField, err error) {
	goodsM = make(map[string]GoodsField)
	var _goodsList []Goods
	if err = utils.GormFind(db, &_goodsList, utils.NewWhereCond("owner_user", ownerUser), utils.NewInCondFromString("uid", uuidList)); err != nil {
		return
	}
	for _, c := range _goodsList {
		goodsM[c.UID] = GoodsField{
			GoodsName: c.Name,
			GoodsTyp:  c.Typ,
		}
	}
	return
}

// goodsM key goods uuid
func SetSurplusByBatch(db *gorm.DB, batchUUID string) (goodsM map[string]*SurplusField, err error) {
	goodsM = make(map[string]*SurplusField)
	var _bgs []BatchGoods
	if err = utils.GormFind(db, &_bgs, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bg := range _bgs {
		goodsM[_bg.GoodsUUID] = &SurplusField{
			Weight: _bg.Weight,
			Mount:  _bg.Mount,
		}
	}
	var _bos []BatchOrderGoods
	if err = utils.GormFind(db, &_bos, utils.NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bo := range _bos {
		if _, ok := goodsM[_bo.GoodsUUID]; ok {
			goodsM[_bo.GoodsUUID] = &SurplusField{
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

func (s *SurplusField) Set() {
	if s.Mount != 0 {
		s.Surplus = fmt.Sprintf("%d", s.Mount)
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.Weight)
	return
}
