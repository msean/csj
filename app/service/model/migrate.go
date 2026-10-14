package model

import (
	"app/global"
)

func Migrate() (err error) {
	global.Global.DB.AutoMigrate(
		&User{},
		&Customer{},
		&GoodsCategory{},
		&Goods{},
		&Batch{},
		&BatchGoods{},
		&BatchOrderHistory{},
		&BatchSerialCounter{},
		&Sms{},
		&MessageCenter{},
		&ShareRecord{},
		&AppVersion{},
	)

	return
}
