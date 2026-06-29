package response

type (
	BatchListItem struct {
		UID      string `json:"uuid"`
		Time     string `json:"time"`
		Status   int    `json:"status"`
		Title    string `json:"title"`
		SerialID int    `json:"serialID"`
	}
	PrecreateRsp struct {
		SerialID int `json:"serialID"`
	}
	BatchOrderGoodsOrderRsp struct {
		GoodsUUID string  `json:"goodsUUID"`
		Price     float64 `json:"price"`
	}
	BatchGoodsGroupItem struct {
		SellTime     string `json:"sellTime"`
		CustomerName string `json:"customerName"`
		SellPrice    string `json:"sellPrice"`
		SellAmount   string `json:"sellAmount"`
		SellTotall   string `json:"sellTotal"`
		GoodType     int    `json:"goodType"`
	}
	BatchGoodsGroupRsp struct {
		Items   []*BatchGoodsGroupItem `json:"items"`
		Profit  string                 `json:"profit"` // 利润
		Surplus string                 `json:"surplus"`
	}
)
