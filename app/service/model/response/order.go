package response

type (
	ShareDailyOrderItem struct {
		GoodsUUID string  `json:"goodsUUID"`
		GoodsName string  `json:"goodsName"` // from GoodsFeild
		Price     float64 `json:"price"`
		Weight    float64 `json:"weight"`
		Mount     int     `json:"mount"`
		Total     float64 `json:"total"`
	}
	ShareDailyOrderRsp struct {
		TotalAmount         float64                `json:"totalAmount"`
		CreditAmount        float64                `json:"creditAmount"`
		GoodsList           []*ShareDailyOrderItem `json:"goodsList"`
		TotalPreviousCredit float64                `json:"totalPreviousCredit"` // sum CreditAmount before today
	}
)
