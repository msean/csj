package request

import "app/utils"

type (
	BatchOrderGoodsParam struct {
		GoodsUUID string  `json:"goodsUUID"`
		Price     float64 `json:"price"`
		Weight    float64 `json:"weight"`
		Mount     int32   `json:"mount"`
		SerialNo  string  `json:"serialNo"`
	}
	CreateTempBatchOrderParam struct {
		BatchUUID    string                 `json:"batchUUID"`
		CustomerUUID string                 `json:"customerUUID"`
		GoodsList    []BatchOrderGoodsParam `json:"goodsList"`
	}
	BatchOrderDetailParam struct {
		UUID string `json:"uuid"`
	}
	CreateBatchOrderParam struct {
		BatchUUID    string                 `json:"batchUUID"`
		CustomerUUID string                 `json:"customerUUID"`
		GoodsList    []BatchOrderGoodsParam `json:"goodsList"`
		FPayAmount   float64                `json:"payAmount"` // 总计
		PayType      int32                  `json:"payType"`   // 支付方式
	}
	UpdateBatchOrderParam struct {
		BatchOrderUUID string                 `json:"uuid"`
		CustomerUUID   string                 `json:"customerUUID"`
		GoodsList      []BatchOrderGoodsParam `json:"goodsList"`
	}
	UpdateBatchOrderStatusParam struct {
		BatchOrderUUID string `json:"uuid"`
		Status         int32  `json:"status"`
	}
	GoodsLatestOrderParam struct {
		GoodsUUIDList []string `json:"goodsUUIDList"`
	}
	ListBatchOrderParam struct {
		utils.LimitCond
		Status    int32  `json:"status"`
		UserUUID  string `json:"userUUID"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
	}
	ShareBatchOrderrParam struct {
		UUID string `json:"uuid"`
	}
)
