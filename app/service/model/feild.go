package model

import "gorm.io/gorm"

type (
	CustomerFeild struct {
		CustomerName  string `gorm:"-"  json:"customerName"`
		CustomerPhone string `gorm:"-"  json:"customerPhone"`
	}
	GoodsFeild struct {
		GoodsName string `gorm:"-" json:"Name"`
		GoodsTyp  int32  `gorm:"-" json:"Typ"`
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
