package response

// 赊欠人统计
type (
	CreditUserStat struct {
		UserUUID    string  `json:"userUUID"`
		UserName    string  `json:"userName"`    // 客户名称
		TotalCredit float64 `json:"totalCredit"` // 总赊欠金额
		TodayCredit float64 `json:"todayCredit"` // 今日赊欠金额
		LongestDays int     `json:"longestDays"` // 最长赊欠天数
		OrderCount  int64   `json:"orderCount"`  // 赊欠订单数
	}

	// 赊欠汇总
	CreditSummary struct {
		TotalCreditAmount float64 `json:"totalCreditAmount"` // 总赊欠金额
		TotalCreditUsers  int64   `json:"totalCreditUsers"`  // 赊欠人数
		TotalCreditOrders int64   `json:"totalCreditOrders"` // 赊欠总单数
		TodayCreditAmount float64 `json:"todayCreditAmount"` // 今日赊欠总额
	}

	// 赊欠列表响应
	CreditListResponse struct {
		List    []*CreditUserStat `json:"list"`
		Summary CreditSummary     `json:"summary"`
	}
)
