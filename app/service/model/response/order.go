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
		Title                string                 `json:"title"`                // 抬头
		CustomerName         string                 `json:"customerName"`         // 客户名称
		Date                 string                 `json:"date"`                 // 日期
		Contact              string                 `json:"contact"`              // 联系方式 放底部左侧
		TotalCreditAmount    float64                `json:"totalCreditAmount"`    // 总赊欠（前欠+今日赊欠）
		TodayCreditAmount    float64                `json:"todayCreditAmount"`    // 今日赊欠
		PreviousCreditAmount float64                `json:"previousCreditAmount"` // 前欠（今日之前的总赊欠）
		GoodsList            []*ShareDailyOrderItem `json:"goodsList"`            // 货品列表
	}
)
