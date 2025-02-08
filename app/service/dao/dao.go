package dao

var (
	User            = NewUserDao()
	Customer        = NewCustomerDao()
	Goods           = NewGoodsDao()
	GoodsCategory   = NewGoodsCategoryDao()
	Batch           = NewBatchDao()
	BatchGoods      = NewBatchGoods()
	BatchOrder      = NewBatchOrderDao()
	BatchOrderGoods = NewBatchOrderGoodsDao()
	BatchOrderPay   = NewBatchOrderPayDao()
)
