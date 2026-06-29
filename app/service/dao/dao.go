package dao

var (
	GoodsDao      = newGoodsDao()
	CustomerDao   = newCustomerDao()
	BatchDao      = newbatchDao()
	BatchGoodsDao = newbatchGoodsDao()
	OrderDao      = newShardOrderDao()
	UserDao       = newUserDao()
	PaymentDao    = newPaymentDao()
)
