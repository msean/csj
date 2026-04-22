package constant

const (
	AdRechargeCreate      = 1 // 创建支付
	AdRechargePaid        = 2 // 完成支付
	AdRechargeTimeout     = 3 // 超时
	AdRechargePaidTimeout = 4 // 超时但完成
	AdRechargeCancel      = 5 // 取消订单
)

const (
	OrderLeftPaid = 10 // 留给用户付款的剩余时间
	OrderMatchAgo = 15 // 匹配交易当前之前的时间
)
