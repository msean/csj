package request

type (
	// QuickPayListReq 快速还款列表请求
	QuickPayListReq struct {
		Page      int `json:"page"`
		PageCount int `json:"pageCount"`
	}

	// QuickPayReq 快捷还款请求
	QuickPayReq struct {
		CustomerUUID string  `json:"customerUUID" binding:"required"`
		Amount       float64 `json:"amount" binding:"required"`
		PayType      int32   `json:"payType"`
		Remark       string  `json:"remark"`
	}

	// OrderPayReq 针对订单还款请求
	OrderPayReq struct {
		OrderUUID string  `json:"orderUUID" binding:"required"`
		Amount    float64 `json:"amount" binding:"required"`
		PayType   int32   `json:"payType"`
		Remark    string  `json:"remark"`
	}

	// PayHistoryReq 还款历史请求
	PayHistoryReq struct {
		CustomerUUID string `json:"customerUUID" binding:"required"`
		Page         int    `json:"page"`
		PageCount    int    `json:"pageCount"`
	}

	// PayDetailReq 还款详情请求
	PayDetailReq struct {
		PayUUID string `json:"payUUID" binding:"required"`
	}

	// RevokePayReq 撤销还款请求
	RevokePayReq struct {
		PayUUID string `json:"payUUID" binding:"required"`
		Reason  string `json:"reason"`
	}

	// MessageListReq 消息列表请求
	MessageListReq struct {
		Page      int `json:"page"`
		PageCount int `json:"pageCount"`
	}
)
