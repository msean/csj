package model

import (
	"app/global"
)

func Migrate() (err error) {
	return global.Global.DB.AutoMigrate(
		&User{},
		&QuickPayAudit{},
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
		&SupportConfig{},
		&OfficialMessageRead{},
		&ShareRecord{},
		&AppVersion{},
		&MemberConfig{},
		&RechargeOrder{},
	)
}
