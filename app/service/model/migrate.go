package model

import (
	"app/global"
)

func Migrate() {
	if global.GlobalRunTime.Migrate() {
		global.GlobalRunTime.DB.AutoMigrate(
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
		)
	}
}
