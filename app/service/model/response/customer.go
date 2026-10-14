package response

import "app/service/model"

type (
	CustomerListItem struct {
		model.Customer
		LatestBillDate int `json:"latestBillDate"`
	}
	CustomerListRsp struct {
		TotalCustomer int64              `json:"totalCustomer"`
		TotalCredit   float64            `json:"totalCredit"`
		Customers     []CustomerListItem `json:"customers"`
	}
	CustomerDetailRsp struct {
		Name string `json:"name"`
		// CreditDays      string  `json:"creditDays"` // 赊欠天数
		// UnOrderDays     string  `json:"unOrderDays"`
		Credit          float64 `json:"credit"`          // 赊欠总金额
		OrderAmount     float64 `json:"orderAmount"`     // 拿货总金额
		OrderCount      float64 `json:"orderCount"`      // 拿货总订单数
		LatestPayDate   string  `json:"latestPayDate"`   // 最近一次还款记录
		LatestPayAmount float64 `json:"latestPayAmount"` // 最近一次还款金额
	}
)
