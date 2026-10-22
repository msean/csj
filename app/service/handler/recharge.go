package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func rechargeRouter(g *gin.RouterGroup) {
	// 需鉴权的接口
	authGroup := g.Group("/recharge", middleware.AuthMiddleware())
	{
		authGroup.POST("/create", RechargeCreate)
		authGroup.POST("/status", RechargeStatus)
		authGroup.POST("/list", RechargeList)
	}
	// 微信支付回调（无需鉴权）
	g.POST("/recharge/notify", RechargeNotify)
}

// RechargeCreate 创建会员充值订单
func RechargeCreate(c *gin.Context) {
	var req request.CreateRechargeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	rsp, err := logic.NewRechargeLogic(c).CreateRechargeOrder(req.MemberType)
	common.Response(c, err, rsp)
}

// RechargeNotify 微信支付回调通知
func RechargeNotify(c *gin.Context) {
	err := logic.NewRechargeLogic(c).WechatNotify()
	if err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": "SUCCESS", "message": "成功"})
}

// RechargeStatus 查询充值订单状态
func RechargeStatus(c *gin.Context) {
	var req request.RechargeStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	rsp, err := logic.NewRechargeLogic(c).QueryRechargeStatus(req.OrderNo)
	common.Response(c, err, rsp)
}

// RechargeList 充值记录列表
func RechargeList(c *gin.Context) {
	var req request.RechargeListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	list, total, err := logic.NewRechargeLogic(c).RechargeList(req.Page, req.PageCount)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, gin.H{"items": list, "total": total})
}
