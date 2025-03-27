package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"
	"app/service/model/response"

	"github.com/gin-gonic/gin"
)

func batchRouter(g *gin.RouterGroup) {
	batchGroup := g.Group("/batch", middleware.AuthMiddleware())
	{
		batchGroup.POST("/surplus", BatchSurplus)
		batchGroup.POST("/detail", BatchDetail)
		batchGroup.POST("/create", BatchCreate)
		batchGroup.POST("/update", BatchUpdate)
		batchGroup.POST("/update/status", BatchUpdateStatus)
	}
	batchGoodsGroup := g.Group("/batch/goods")
	{
		batchGoodsGroup.POST("/update", BatchGoodsUpdate)
		batchGoodsGroup.POST("/detail", BatchGoodsDetail)
	}
}

func BatchSurplus(c *gin.Context) {
	batchGoodsList, err := logic.NewBatchLogic(c).CalSurplus()
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, gin.H{
		"surplus": batchGoodsList,
	})
}

func BatchCreate(c *gin.Context) {
	var param request.CreateBatchParam
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchLogic := logic.NewBatchLogic(c)
	var err error
	err = batchLogic.Create(param)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func BatchUpdate(c *gin.Context) {
	var param request.UpdateBatchParam

	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchLogic := logic.NewBatchLogic(c)
	var err error
	if err = batchLogic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func BatchUpdateStatus(c *gin.Context) {
	var param request.UpdateBatchStatusParam
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic := logic.NewBatchLogic(c)
	if err := batchLogic.UpdateStatus(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func BatchDetail(c *gin.Context) {
	var param request.BatchDetailParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic := logic.NewBatchLogic(c)
	var rsp response.BatchRsp
	rsp, err = batchLogic.Detail(param)
	common.Response(c, err, rsp)
}

func BatchGoodsDetail(c *gin.Context) {
	var param request.FindBatchGoodsParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchGoods response.BatchGoodsRsp
	if batchGoods, err = logic.NewBatchGoodsLogic(c).FromUUID(param.UUIDCompatible); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchGoods)
}

func BatchGoodsUpdate(c *gin.Context) {
	var param request.UpdateBatchGoodsParam
	var err error
	if err = c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchGoods response.BatchGoodsRsp
	if batchGoods, err = logic.NewBatchGoodsLogic(c).Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchGoods)
}
