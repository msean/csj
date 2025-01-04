package request

type (
	SaveBatchGoods struct {
		UUID   string  `json:"goodsUUID"`
		Price  float64 `json:"price"`
		Weight float64 `json:"weight"`
		Mount  int32   `json:"mount"`
	}
	CreateBatchParam struct {
		StroageTime int64            `json:"storageTime"`
		Goods       []SaveBatchGoods `json:"goodsList"`
	}
	BatchDetailParam struct {
		UUID string `json:"uuid"`
		Date string `json:"date"`
	}
	UpdateBatchParam struct {
		StroageTime int64            `json:"storageTime"`
		UUID        string           `json:"uuid"`
		Goods       []SaveBatchGoods `json:"goodsList"`
	}
	UpdateBatchStatusParam struct {
		UUID   string `json:"uuid"`
		Status int    `json:"status"`
	}
	FindBatchGoodsParam struct {
		UUID string `json:"uuid"`
	}
	UpdateBatchGoodsParam struct {
		SaveBatchGoods
	}
)
