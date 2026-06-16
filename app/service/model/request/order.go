package request

import "app/pkg/utils"

type (
	ShareDailyOrderReq struct {
		CustomerUUID string `json:"customerUUID" binding:"required"` // UserUUID of the customer
	}
	CreditListReq struct {
		utils.LimitCond
	}
)
