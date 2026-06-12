package response

import "app/service/model"

type (
	UserProfileRsp struct {
		model.User
		MonthSales   string `json:"monthSales"`   // 月销售额
		CreditAmount string `json:"customerDebt"` // 赊欠总金额
	}
)
