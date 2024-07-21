package model

import (
	"fmt"

	"gorm.io/gorm"
)

type (
	SurplusFeild struct {
		Weight  float32 `gorm:"-" json:"-"`
		Mount   int32   `gorm:"-" json:"-"`
		Surplus string  `gorm:"-" json:"surplus"`
	}
	CustomerFeild struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsFeild struct {
		GoodsName string `gorm:"-" json:"name"`
		GoodsTyp  int32  `gorm:"-" json:"type"`
	}
	PayFeild struct {
		PayFee  float32 `gorm:"-"  json:"payFee"`
		PayType int32   `gorm:"-"  json:"payType"`
		PaidFee float32 `gorm:"-"  json:"paidFee"`
	}
)

func CustomerFeildSet(db *gorm.DB, uuid, ownerUser string) (c CustomerFeild, err error) {
	var customer Customer
	if err = Find(db, &customer, NewWhereCond("owner_user", ownerUser), WhereUIDCond(uuid)); err != nil {
		return
	}
	c.CustomerName = customer.Name
	c.CustomerPhone = customer.Phone
	return
}

func BatchCustomerFeildSet(db *gorm.DB, uuidList []string, ownerUser string) (customerM map[string]CustomerFeild, err error) {
	customerM = make(map[string]CustomerFeild)
	var customers []Customer
	if err = Find(db, &customers, NewWhereCond("owner_user", ownerUser), NewInCondFromString("uid", uuidList)); err != nil {
		return
	}
	for _, c := range customers {
		customerM[c.UID] = CustomerFeild{
			CustomerName:  c.Name,
			CustomerPhone: c.Phone,
		}
	}
	return
}

func GoodsFeildSet(db *gorm.DB, uuid, ownerUser string) (c GoodsFeild, err error) {
	var _goods Goods
	if err = Find(db, &_goods, NewWhereCond("owner_user", ownerUser), WhereUIDCond(uuid)); err != nil {
		return
	}
	c = GoodsFeild{
		GoodsName: _goods.Name,
		GoodsTyp:  _goods.Typ,
	}
	return
}

func BatchGoodsFeildSet(db *gorm.DB, uuidList []string, ownerUser string) (goodsM map[string]GoodsFeild, err error) {
	goodsM = make(map[string]GoodsFeild)
	var _goodsList []Goods
	if err = Find(db, &_goodsList, NewWhereCond("owner_user", ownerUser), NewInCondFromString("uid", uuidList)); err != nil {
		return
	}
	for _, c := range _goodsList {
		goodsM[c.UID] = GoodsFeild{
			GoodsName: c.Name,
			GoodsTyp:  c.Typ,
		}
	}
	return
}

// goodsM key goods uuid
func SetSurplusByBatch(db *gorm.DB, batchUUID string) (goodsM map[string]*SurplusFeild, err error) {
	goodsM = make(map[string]*SurplusFeild)
	var _bgs []BatchGoods
	if err = Find(db, &_bgs, NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bg := range _bgs {
		goodsM[_bg.GoodsUUID] = &SurplusFeild{
			Weight: _bg.Weight,
			Mount:  _bg.Mount,
		}
	}
	var _bos []BatchOrderGoods
	if err = Find(db, &_bos, NewWhereCond("batch_uuid", batchUUID)); err != nil {
		return
	}
	for _, _bo := range _bos {
		if _, ok := goodsM[_bo.GoodsUUID]; ok {
			goodsM[_bo.GoodsUUID] = &SurplusFeild{
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

func (s *SurplusFeild) Set() {
	if s.Mount != 0 {
		s.Surplus = fmt.Sprintf("%d", s.Mount)
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.Weight)
	return
}
