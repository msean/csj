package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"

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
	var batch model.Batch
	var err error
	batch, err = batchLogic.Create(param)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic.SetGoodsFeild(batch)
	common.Response(c, nil, batch)
}

func BatchUpdate(c *gin.Context) {
	var param request.UpdateBatchParam

	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchLogic := logic.NewBatchLogic(c)
	var batch model.Batch
	var err error
	if batch, err = batchLogic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batch)
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
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic := logic.NewBatchLogic(c)
	var err error
	var batch model.Batch
	if param.UUID != "" {
		batch, err = batchLogic.FromUUID(param.UUID, true)
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, batch)
		return
	}

	if param.Date != "" {
		batch, err = batchLogic.FromDate(param.Date)
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, batchLogic)
		return
	}

	batch, err = batchLogic.FromLatest()
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchLogic)
}

func BatchGoodsDetail(c *gin.Context) {
	var param request.FindBatchGoodsParam

	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	var batchGoods model.BatchGoods
	var err error
	if batchGoods, err = logic.NewBatchGoodsLogic(c).FromUUID(param.UUID); err != nil {
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

	var batchGoods model.BatchGoods
	if batchGoods, err = logic.NewBatchGoodsLogic(c).Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchGoods)
}
