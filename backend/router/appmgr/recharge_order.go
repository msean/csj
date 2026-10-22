package appmgr

import (
	"github.com/gin-gonic/gin"
)

type RechargeOrderRouter struct{}

// InitRechargeOrderRouter 初始化 充值订单 路由信息
func (s *RechargeOrderRouter) InitRechargeOrderRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	rechargeOrderRouter := Router.Group("rechargeOrder")
	{
		rechargeOrderRouter.GET("findRechargeOrder", rechargeOrderApi.FindRechargeOrder)       // 根据uid获取充值订单
		rechargeOrderRouter.GET("getRechargeOrderList", rechargeOrderApi.GetRechargeOrderList) // 获取充值订单列表
	}
}
