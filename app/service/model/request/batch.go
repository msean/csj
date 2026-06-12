package request

import (
	"app/pkg/utils"
	"app/service/model"
)

type (
	BatchOrderDetailReq struct {
		UUID string `json:"uuid"`
	}
	BatchListReq struct {
		utils.LimitCond
		Status    int    `json:"status"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}
	BatchDetailReq struct {
		UUID     string `json:"uuid"`
		Date     string `json:"date"`
		ShowType int    `json:"showType"` // 0 批次列表 1 结算管理
	}
	BatchGoodsListReq struct {
		BatchUUID string `json:"batchUUID"`
		GoodsUUID string `json:"goodsUUID"`
	}
	BatchOrderGoodsLatest struct {
		GoodsUUIDList []string `json:"goodsUUIDList"`
	}
	BatchOrderListReq struct {
		utils.LimitCond
		Status    int32  `json:"status"`
		UserUUID  string `json:"userUUID"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}
	OrderReq struct {
		model.BatchOrder
		FPayAmount float64 `json:"payAmount"` // 付款金额
		// FCreditAmount float64 `json:"creditAmount"`
		PayType int32 `json:"payType"` // 支付方式
	}
)
