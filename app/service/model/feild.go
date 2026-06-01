package model

import (
	"app/global"
	"app/service/common"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type (
	SurplusFeild struct {
		Weight          float64 `gorm:"-" json:"-"`
		Mount           int     `gorm:"-" json:"-"`
		SellMount       int     `gorm:"-" json:"-"`
		SellWeight      float64 `gorm:"-" json:"-"`
		Surplus         string  `gorm:"-" json:"surplus"` // 剩余
		Type            int     `gorm:"-" json:"-"`
		SellAmount      string  `gorm:"-" json:"sellAmount"`      // 销量 件数/均价
		SellAvgPrice    string  `gorm:"-" json:"sellAvgPrice"`    // 销售均价
		SellTotalAmount string  `gorm:"-" json:"sellTotalAmount"` // 销售货款总金额
	}
	CustomerFeild struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsFeild struct {
		GoodsName   string  `gorm:"-" json:"name"`
		GoodsTyp    int32   `gorm:"-" json:"type"`
		GoodsWeight float64 `gorm:"-" json:"goodsWeight"` // 定装是多少斤
	}
	PayFeild struct {
		PayFee  string `gorm:"-"  json:"payFee"`
		PayType int32  `gorm:"-"  json:"payType"`
		PaidFee string `gorm:"-"  json:"paidFee"`
	}
	BatchStatFeild struct {
		StatMount     int `json:"statMount"`     // 总件数
		StatWeight    int `json:"statWeight"`    // 总重量
		StatInventory int `json:"StatInventory"` // 库存
		// StatInventory int `json:"StatInventory"` // 库存
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
		GoodsName:   _goods.Name,
		GoodsTyp:    _goods.Typ,
		GoodsWeight: _goods.Weight,
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
			GoodsName:   c.Name,
			GoodsTyp:    c.Typ,
			GoodsWeight: c.Weight,
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
			Weight:     _bg.Weight,
			Mount:      _bg.Mount,
			Type:       _bg.GoodType,
			SellAmount: "0",
		}
	}
	global.Global.Logger.Debug("SetSurplusByBatch _bgs", zap.Any("goodsM", goodsM))
	var _bos []BatchOrderGoods
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
			)
			goodsM[_bo.GoodsUUID] = &SurplusFeild{
				Weight:     goodsM[_bo.GoodsUUID].Weight - _bo.Weight,
				Mount:      goodsM[_bo.GoodsUUID].Mount - _bo.Mount,
				Type:       goodsM[_bo.GoodsUUID].Type,
				SellWeight: goodsM[_bo.GoodsUUID].SellWeight + _bo.Weight,
				SellMount:  goodsM[_bo.GoodsUUID].Mount + _bo.Mount,
			}
		}
	}

	global.Global.Logger.Info("SetSurplusByBatch", zap.Any("goodsM", goodsM))
	for _, surplus := range goodsM {
		global.Global.Logger.Info("SetSurplusByBatch",
			zap.Any("surplus Weight", surplus.Weight),
			zap.Any("surplus Mount", surplus.Mount),
			zap.Any("surplus Surplus", surplus.Surplus),
			zap.Any("surplus Type", surplus.Type),
		)
		surplus.Set()
	}
	return
}

func (s *SurplusFeild) Set() {
	if s.Type == common.GoodsTypeFix {
		s.Surplus = fmt.Sprintf("%d", s.Mount)
		s.SellAmount = fmt.Sprintf("%d", s.SellMount)
		return
	}
	s.Surplus = fmt.Sprintf("%.2f", s.Weight)
	s.SellAmount = fmt.Sprintf("%.2f", s.SellWeight)
}
