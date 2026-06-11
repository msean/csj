package request

import "app/pkg/utils"

type BatchListReq struct {
	utils.LimitCond
	Status    int    `json:"status"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type BatchDetailReq struct {
	UUID     string `json:"uuid"`
	Date     string `json:"date"`
	ShowType int    `json:"showType"` // 0 批次列表 1 结算管理
}

type BatchGoodsListReq struct {
	BatchUUID string `json:"batchUUID"`
	GoodsUUID string `json:"goodsUUID"`
}
