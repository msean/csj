package response

import (
	"app/pkg/wechat"
	"app/service/model"
	"time"
)

type (
	// CreateRechargeRsp 创建充值订单响应
	CreateRechargeRsp struct {
		OrderNo string             `json:"orderNo"`
		PayInfo *wechat.AppPayParams `json:"payInfo"`
	}

	// RechargeStatusRsp 充值状态响应
	RechargeStatusRsp struct {
		OrderNo     string    `json:"orderNo"`
		PayStatus   int       `json:"payStatus"`
		PayTime     *time.Time `json:"payTime"`
		VipDeadline time.Time `json:"vipDeadline"`
	}

	// RechargeListItem 充值记录项
	RechargeListItem struct {
		model.RechargeOrder
	}
)
