package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func paymentRouter(g *gin.RouterGroup) {
	paymentGroup := g.Group("/payment", middleware.AuthMiddleware())
	{
		// 快速还款列表
		paymentGroup.POST("/quick/list", QuickPayList)
		// 快捷还款
		paymentGroup.POST("/quick/pay", QuickPay)
		// 针对订单还款
		paymentGroup.POST("/order/pay", OrderPay)
		// 还款历史
		paymentGroup.POST("/history", PayHistory)
		// 还款详情
		paymentGroup.POST("/detail", PayDetail)
		// 撤销还款
		paymentGroup.POST("/revoke", RevokePay)
		// 还款分享
		paymentGroup.POST("/share/generate_html", PaymentShareGenerateHTML)
	}
}

// QuickPayList 快速还款列表
func QuickPayList(c *gin.Context) {
	var req request.QuickPayListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	if req.PageCount <= 0 {
		req.PageCount = 10
	}

	items, total, totalDebt, err := logic.NewPaymentLogic(c).QuickPayList(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, gin.H{
		"items":         items,
		"totalCustomer": total,
		"totalCredit":   totalDebt,
	})
}

// QuickPay 快捷还款
func QuickPay(c *gin.Context) {
	var req request.QuickPayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewPaymentLogic(c).QuickPay(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, rsp)
}

// OrderPay 针对订单还款
func OrderPay(c *gin.Context) {
	var req request.OrderPayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewPaymentLogic(c).OrderPay(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, rsp)
}

// PayHistory 还款历史
func PayHistory(c *gin.Context) {
	var req request.PayHistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	if req.PageCount <= 0 {
		req.PageCount = 10
	}

	items, total, err := logic.NewPaymentLogic(c).PayHistory(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, gin.H{
		"items": items,
		"total": total,
	})
}

// PayDetail 还款详情
func PayDetail(c *gin.Context) {
	var req request.PayDetailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	rsp, err := logic.NewPaymentLogic(c).PayDetail(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, rsp)
}

// RevokePay 撤销还款
func RevokePay(c *gin.Context) {
	var req request.RevokePayReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	err := logic.NewPaymentLogic(c).RevokePay(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, nil)
}

// ShareGenerateHTML 生成今日账单HTML，直接返回HTML内容
func PaymentShareGenerateHTML(c *gin.Context) {
	var req request.PaymentShareReq
	if err := c.ShouldBind(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	html, err := logic.NewPaymentLogic(c).GenerateShareHTML(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
