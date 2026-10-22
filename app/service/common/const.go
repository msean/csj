package common

const (
	GoodsTypeFix  = 1 //定装
	GoodsTypeBulk = 2 // 散装
)

const (
	BatchOrderTemp     int = 1   // 记账单
	BatchOrderedCredit int = 2   // 赊欠单 废弃
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

const (
	BatchStatusOnSellering = 1 // 在售
	BatchStatusOffShelf    = 2 // 售罄
	BatchStatusSettled     = 3 // 已结算
)

const (
	CashCustomerName = "现金客户" // 现金客户名称，多处使用，修改时改这里即可
)

const (
	MemberTypeMonth   = 1 // 月
	MemberTypeQuarter = 2 // 季度
	MemberTypeYear    = 3 // 年
)

const (
	PayDetailTypeDebt  = 1 // 还原始客户债务(PaidDebt)
	PayDetailTypeOrder = 2 // 还订单赊欠
)

var ExCludeTempBatchOrder = []int{
	BatchOrderedCredit, BatchOrderFinish, BatchOrderCancel, BatchOrderRefund, BatchOrderReTurn,
}

var ValidOrder = []int{
	BatchOrderTemp, BatchOrderFinish,
}

var UnFinalStatusOrder = []int{
	BatchOrderTemp,
}
