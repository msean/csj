package common

const (
	GoodsTypeFix  = 1 //定装
	GoodsTypeBulk = 2 // 散装
)

const (
	BatchOrderTemp     int = 1   // 记账单
	BatchOrderedCredit int = 2   // 赊欠单
	BatchOrderFinish   int = 3   // 已结算
	BatchOrderCancel   int = 100 //作废
	BatchOrderRefund   int = 101 // 退款
	BatchOrderReTurn   int = 102 // 退货

	BatchOrderUnshare = 1 // 未分享
	BatchOrderShared  = 2 // 已分享

	PayTypeWx   = 1 // 微信
	PayTypeZFB  = 2 // 支付宝
	PayTypeBank = 3 // 银行
	PayTypeCash = 4 // 现金
)

var ExCludeTempBatchOrder = []int{
	BatchOrderedCredit, BatchOrderFinish, BatchOrderCancel, BatchOrderRefund, BatchOrderReTurn,
}

var FinalBatchOrder = []int{
	BatchOrderReTurn, BatchOrderCancel, BatchOrderRefund,
}

var ValidOrder = []int{
	BatchOrderTemp, BatchOrderedCredit, BatchOrderFinish,
}
