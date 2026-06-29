package model

import (
	"app/global"
)

func Migrate() {
	global.Global.DB.AutoMigrate(
		&User{},
		&Customer{},
		&GoodsCategory{},
		&Goods{},
		&Batch{},
		&BatchGoods{},
		&BatchOrder{},
		&BatchOrderGoods{},
		&BatchOrderPay{},
		&BatchOrderHistory{},
		&BatchSerialCounter{},
		&Sms{},
		&MessageCenter{},
	)
}
