package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
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
		batchGroup.POST("/preCreate", BatchPrecreate)
		batchGroup.POST("/list", BatchList)
		batchGroup.POST("/update/status", BatchUpdateStatus)
	}
	batchGoodsGroup := g.Group("/batch/goods", middleware.AuthMiddleware())
	{
		batchGoodsGroup.POST("/update", BatchGoodsUpdate)
		batchGoodsGroup.POST("/detail", BatchGoodsDetail)
		batchGoodsGroup.POST("/list", OrderGoodsList)
	}
}

func BatchSurplus(c *gin.Context) {
	var batchGoodsList []logic.BatchGoodsLogic
	var err error
	batchGoodsList, err = logic.NewBatchLogic(c).CalSurplus()
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
	var err error
	if err = c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	err = batchLogic.Update(true)
	common.Response(c, err, batchLogic)
}

func BatchPrecreate(c *gin.Context) {
	batchLogic := logic.NewBatchLogic(c)
	var err error
	if err = c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	var rsp response.PrecreateRsp
	rsp, err = batchLogic.Precreate()
	common.Response(c, err, rsp)
}

func BatchList(c *gin.Context) {
	var payload request.BatchListReq
	var err error
	if err = c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchLogic := logic.NewBatchLogic(c)
	var list []response.BatchListItem
	list, err = batchLogic.List(payload)
	common.Response(c, err, list)
}

func BatchUpdateStatus(c *gin.Context) {
	batchLogic := logic.NewBatchLogic(c)
	var err error
	if err = c.ShouldBind(&batchLogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	err = batchLogic.Update(false)
	common.Response(c, err, nil)
}

func BatchDetail(c *gin.Context) {
	var payLoad request.BatchDetailReq
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

	// if payLoad.Date != "" {
	// 	err = batchLogic.FromDate(payLoad.Date)
	// 	if err != nil {
	// 		common.Response(c, err, nil)
	// 		return
	// 	}
	// 	common.Response(c, nil, batchLogic)
	// 	return
	// }

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
	var err error

	if err := c.ShouldBind(&payLoad); err != nil {
		common.Response(c, err, nil)
		return
	}

	batchGoods := logic.NewBatchGoodsLogic(c)
	err = batchGoods.FromUUID(payLoad.UUID)
	common.Response(c, err, batchGoods)
}

func BatchGoodsUpdate(c *gin.Context) {
	batchGoods := logic.NewBatchGoodsLogic(c)
	var err error
	if err = c.ShouldBind(&batchGoods); err != nil {
		common.Response(c, err, nil)
		return
	}

	err = batchGoods.Update()
	common.Response(c, err, batchGoods)
}
