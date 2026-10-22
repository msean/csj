package request

import "github.com/msean/csj/backend/model/common/request"

// RechargeOrderSearch 充值订单搜索条件
type RechargeOrderSearch struct {
	UserUUID  string `json:"userUUID" form:"userUUID"`
	OrderNo   string `json:"orderNo" form:"orderNo"`
	PayStatus int    `json:"payStatus" form:"payStatus"`
	request.PageInfo
}
