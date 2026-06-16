package request

import "app/pkg/utils"

type (
	CustomerListReq struct {
		utils.LimitCond
		SearchKey string `json:"searchName"`
	}
)
