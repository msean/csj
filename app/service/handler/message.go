package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func messageRouter(g *gin.RouterGroup) {
	messageGroup := g.Group("/message", middleware.AuthMiddleware())
	{
		// 消息列表
		messageGroup.POST("/list", MessageList)
		// 标记消息已读
		messageGroup.POST("/read", MessageRead)
	}
}

// MessageList 消息列表
func MessageList(c *gin.Context) {
	var req request.MessageListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}
	if req.PageCount <= 0 {
		req.PageCount = 20
	}

	items, summary, err := logic.NewPaymentLogic(c).MessageList(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, gin.H{
		"items":   items,
		"summary": summary,
	})
}

// MessageRead 标记消息已读
func MessageRead(c *gin.Context) {
	var req request.MessageReadReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Response(c, err, nil)
		return
	}

	err := logic.NewPaymentLogic(c).MessageRead(req)
	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, nil)
}
