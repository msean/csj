package model

import (
	"app/global"
)

func Migrate() {
	if global.Global.Migrate() {
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
			&Sms{},
		)
	}
}
