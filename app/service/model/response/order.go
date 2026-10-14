package response

type (
	ShareDailyOrderItem struct {
		GoodsUUID string  `json:"goodsUUID"`
		GoodsName string  `json:"goodsName"` // from GoodsFeild
		Price     float64 `json:"price"`
		Weight    float64 `json:"weight"`
		Mount     int     `json:"mount"`
		Total     float64 `json:"total"`
		GoodType  int     `json:"goodType"` // 1-定装 2-散装
	}
	// ShareDailyOrderGroup 单个订单的货品分组（用于分享时按订单分开渲染）
	ShareDailyOrderGroup struct {
		OrderUUID   string                 `json:"orderUUID"`
		DateStr     string                 `json:"dateStr"`     // 开单时间
		GoodsList   []*ShareDailyOrderItem `json:"goodsList"`
		TricycleFee float64                `json:"tricycleFee"`
		TotalAmount float64                `json:"totalAmount"` // 该订单赊欠金额
		TotalMount  int                    `json:"totalMount"`
		TotalWeight float64                `json:"totalWeight"`
	}
	ShareDailyOrderRsp struct {
		Title                string                 `json:"title"`                // 抬头
		CustomerName         string                 `json:"customerName"`         // 客户名称
		Date                 string                 `json:"date"`                 // 日期
		Contact              string                 `json:"contact"`              // 联系方式 放底部左侧
		TotalCreditAmount    float64                `json:"totalCreditAmount"`    // 总赊欠（前欠+今日赊欠）
		TodayCreditAmount    float64                `json:"todayCreditAmount"`    // 今日赊欠
		PreviousCreditAmount float64                `json:"previousCreditAmount"` // 前欠（今日之前的总赊欠）
		TricycleFee          float64                `json:"tricycleFee"`          // 三轮车费用（多订单累加）
		OrderGroups          []ShareDailyOrderGroup `json:"orderGroups"`          // 按订单分组
	}
)
