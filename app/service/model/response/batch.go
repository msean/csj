package response

type (
	BatchListItem struct {
		UID    string `json:"uuid"`
		Time   string `json:"time"`
		Status int    `json:"status"`
		Title  string `json:"title"`
	}
)

type PrecreateRsp struct {
	SerialID int `json:"serialID"`
}

type BatchGoodsGroupItem struct {
	SellTime     string `json:"sellTime"`
	CustomerName string `json:"customerName"`
	SellPrice    string `json:"sellPrice"`
	SellAmount   string `json:"sellAmount"`
	SellTotall   string `json:"sellTotal"`
	GoodType     int    `json:"goodType"`
}

type BatchGoodsGroupRsp struct {
	Items []BatchGoodsGroupItem `json:"items"`
}
