package dao

import "app/pkg/utils"

var (
	GoodsDao      = newGoodsDao()
	CustomerDao   = newCustomerDao()
	BatchDao      = newbatchDao()
	BatchGoodsDao = newbatchGoodsDao()
	OrderDao      = newOrderDao()
	UserDao       = newUserDao()
	PaymentDao    = newPaymentDao()
	ShareDao      = newShareDao()
)

func ShardOrderTable(ownerUser string) string {
	return utils.GetShardTableName("batch_orders", ownerUser)
}

func ShardOrderGoodsTable(ownerUser string) string {
	return utils.GetShardTableName("batch_order_goods", ownerUser)
}
func ShardPayTable(ownerUser string) string {
	return utils.GetShardTableName("batch_order_pay", ownerUser)
}
