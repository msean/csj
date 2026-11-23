package response

import "time"

type (
	// QuickPayListItem 快速还款列表项
	QuickPayListItem struct {
		CustomerUUID string  `json:"customerUUID"`
		CustomerName string  `json:"customerName"`
		TotalCredit  float64 `json:"totalCredit"` // 总赊欠金额
	}

	// QuickPayRsp 快捷还款响应
	QuickPayRsp struct {
		PayUUID  string  `json:"payUUID"`
		Discount float64 `json:"discount"` // 优惠
		Amount   float64 `json:"amount"`   // 金额
		NeedPop  bool    `json:"needPop"`  // 是否弹窗
	}

	// OrderRepaymentSummary 订单详情发起还款时的应收金额。
	OrderRepaymentSummary struct {
		PreviousCreditAmount float64 `json:"previousCreditAmount"`
		TotalCreditAmount    float64 `json:"totalCreditAmount"`
	}

	// OrderPayRsp 订单还款响应
	OrderPayRsp struct {
		PayUUID     string  `json:"payUUID"`
		OrderUUID   string  `json:"orderUUID"`
		OrderCredit float64 `json:"orderCredit"` // 还款后订单赊欠金额
	}

	// PayHistoryItem 还款历史项
	PayHistoryItem struct {
		Allocations   []RepaymentAllocation `json:"allocations"`
		PayUUID       string                `json:"payUUID"`
		CustomerUUID  string                `json:"customerUUID"`
		CustomerName  string                `json:"customerName"`
		Amount        float64               `json:"amount"`
		BeforeCredit  float64               `json:"beforeCredit"` // 还款前欠款
		AfterCredit   float64               `json:"afterCredit"`  // 还款后剩余欠款
		PayType       int32                 `json:"payType"`
		Remark        string                `json:"remark"`
		IsRevoked     int                   `json:"isRevoked"`
		RevokedAt     *time.Time            `json:"revokedAt"`
		RevokedReason string                `json:"revokedReason"`
		CreatedAt     time.Time             `json:"createdAt"`
	}

	// PayDetailRsp 还款详情响应
	PayDetailRsp struct {
		Allocations     []RepaymentAllocation `json:"allocations"`
		RecordOwnerName string                `json:"recordOwnerName"`
		PayUUID         string                `json:"payUUID"`
		CustomerUUID    string                `json:"customerUUID"`
		CustomerName    string                `json:"customerName"`
		OrderUUID       string                `json:"orderUUID"` // 订单UUID（快捷还款为空）
		Amount          float64               `json:"amount"`
		BeforeCredit    float64               `json:"beforeCredit"` // 还款前欠款
		AfterCredit     float64               `json:"afterCredit"`  // 还款后剩余欠款
		PayType         int32                 `json:"payType"`
		Remark          string                `json:"remark"`
		PayDetails      string                `json:"payDetails"` // 快捷还款的订单详情JSON
		IsRevoked       int                   `json:"isRevoked"`
		RevokedAt       *time.Time            `json:"revokedAt"`
		RevokedReason   string                `json:"revokedReason"`
		CreatedAt       time.Time             `json:"createdAt"`
		Discount        float64               `json:"discount"`
	}

	// MessageItem 消息项
	MessageItem struct {
		Source       string    `json:"source"`
		Title        string    `json:"title"`
		MessageUUID  string    `json:"messageUUID"`
		Type         int       `json:"type"`
		Event        string    `json:"event"`
		Content      string    `json:"content"`
		CustomerUUID string    `json:"customerUUID"`
		CustomerName string    `json:"customerName"`
		IsRead       int       `json:"isRead"`
		RelatedUUID  string    `json:"relatedUUID"`
		CreatedAt    time.Time `json:"createdAt"`
	}

	// MessageSummary 消息摘要
	MessageSummary struct {
		TotalCount  int `json:"totalCount"`
		UnreadCount int `json:"unreadCount"`
	}
)

// RepaymentAllocation combines the stored payment allocation with current order data.
// Deleted orders retain their ticket number but cannot be opened.
type RepaymentAllocation struct {
	TicketNo         string     `json:"ticketNo,omitempty"`
	Type             int        `json:"type"`
	OrderUUID        string     `json:"orderUUID"`
	Amount           float64    `json:"amount"`
	OrderAvailable   bool       `json:"orderAvailable"`
	OrderTime        *time.Time `json:"orderTime,omitempty"`
	OrderTotalAmount *float64   `json:"orderTotalAmount,omitempty"`
	OrderStatus      *int       `json:"orderStatus,omitempty"`
}
