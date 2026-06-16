package response

type (
	CustomerDetailRsp struct {
		Name string `json:"name"`
		// CreditDays      string  `json:"creditDays"` // 赊欠天数
		// UnOrderDays     string  `json:"unOrderDays"`
		CreditAmount    float64 `json:"creditAmount"`    // 赊欠总金额
		OrderAmount     float64 `json:"orderAmount"`     // 拿货总金额
		OrderCount      float64 `json:"orderCount"`      // 拿货总订单数
		LatestPayDate   string  `json:"latestPayDate"`   // 最近一次还款记录
		LatestPayAmount float64 `json:"latestPayAmount"` // 最近一次还款金额
	}
)
