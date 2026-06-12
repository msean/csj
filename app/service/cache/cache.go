package cache

import (
	"app/global"
)

var GoodsCache *goodsCache
var CustomerCache *customerCache
var SmsCache *smsCache

func Init() {
	GoodsCache = newGoodsCache(global.Global.Redis, global.Global.DB)
	CustomerCache = newCustomerCache(global.Global.Redis, global.Global.DB)
	SmsCache = NewSmsCache(global.Global.Redis)
}
