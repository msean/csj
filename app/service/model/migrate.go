package model

import (
	"app/global"
)

func Migrate() (err error) {
	return global.Global.DB.AutoMigrate(
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
		&OfficialMessage{},
		&OfficialMessageRead{},
		&ShareRecord{},
		&AppVersion{},
		&MemberConfig{},
		&RechargeOrder{},
	)
}
