package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"

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
	var batchGoodsList []logic.BatchGoodsLogic
	var err error
	batchGoodsList, err = logic.NewBatchLogic(c).CalSurplus(common.GetUserUUID(c))
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, gin.H{
		"surplus": batchGoodsList,
	})
}

func BatchCreate(c *gin.Context) {
	batchLogic := logic.NewBatchLogic(c)
	if err := c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	err := batchLogic.Create()
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic.SetGoodsFeild()
	common.Response(c, nil, batchLogic)
}

func BatchUpdate(c *gin.Context) {
	batchLogic := logic.NewBatchLogic(c)
	if err := c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	if err := batchLogic.Update(true); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchLogic)
}

func BatchUpdateStatus(c *gin.Context) {
	batchLogic := logic.NewBatchLogic(c)
	if err := c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	batchLogic.OwnerUser = common.GetUserUUID(c)
	if err := batchLogic.Update(false); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func BatchDetail(c *gin.Context) {
	type PayLoad struct {
		UUID string `json:"uuid"`
		Date string `json:"date"`
	}
	var payLoad PayLoad
	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchLogic := logic.NewBatchLogic(c)
	var err error
	if payLoad.UUID != "" {
		err = batchLogic.FromUUID(payLoad.UUID)
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, batchLogic)
		return
	}

	if payLoad.Date != "" {
		err = batchLogic.FromDate(payLoad.Date)
		if err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, batchLogic)
		return
	}

	err = batchLogic.FromLatest()
	if err != nil {
		common.Response(c, nil, model.Batch{})
		return
	}
	common.Response(c, nil, batchLogic)
}

func BatchGoodsDetail(c *gin.Context) {
	type PayLoad struct {
		UUID string `json:"uuid"`
	}
	var payLoad PayLoad

	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchGoods := logic.NewBatchGoodsLogic(c)
	if err := batchGoods.FromUUID(payLoad.UUID); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchGoods)
}

func BatchGoodsUpdate(c *gin.Context) {
	batchGoods := logic.NewBatchGoodsLogic(c)
	if err := c.ShouldBind(&batchGoods); err != nil {
		common.Response(c, err, nil)
		return
	}

	if err := batchGoods.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, batchGoods)
}
