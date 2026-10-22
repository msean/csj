package request

type (
	// CreateRechargeReq 创建充值订单请求
	CreateRechargeReq struct {
		MemberType int `json:"memberType" binding:"required"` // 会员类型 1月 2季度 3年
	}

	// RechargeStatusReq 查询充值状态请求
	RechargeStatusReq struct {
		OrderNo string `json:"orderNo" binding:"required"`
	}

	// RechargeListReq 充值记录列表请求
	RechargeListReq struct {
		Page      int `json:"page"`
		PageCount int `json:"pageCount"`
	}
)
