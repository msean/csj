package cache

import (
	"app/global"
)

var GoodsCache *goodsCache
var CustomerCache *customerCache

func Init() {
	GoodsCache = newGoodsCache(global.Global.Redis, global.Global.DB)
	CustomerCache = newCustomerCache(global.Global.Redis, global.Global.DB)
}
