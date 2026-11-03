package request

import "app/pkg/utils"

type (
	ShareDailyOrderReq struct {
		CustomerUUID string `json:"customerUUID" binding:"required"` // UserUUID of the customer
		OrderUUID    string `json:"orderUUID"`
	}
	CreditListReq struct {
		utils.LimitCond
		Search string `json:"search"`
	}
)

type CustomerPurchasesReq struct {
	UserUUID   string `json:"userUUID" binding:"required"`
	Page       int    `json:"page"`
	PageCount  int    `json:"pageCount"`
	StartTime  int64  `json:"startTime" binding:"min=0"`
	EndTime    int64  `json:"endTime" binding:"min=0"`
	CreditOnly bool   `json:"creditOnly"`
}
