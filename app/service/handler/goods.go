package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
)

func goodsRouter(g *gin.RouterGroup) {
	group := g.Group("/goods", middleware.AuthMiddleware())
	{
		group.POST("/save", GoodsSave)
		group.POST("/list", GoodsList)
	}
	categoryGroup := g.Group("/goods/category", middleware.AuthMiddleware())
	{
		categoryGroup.POST("/save", GoodsCategorySave)
		categoryGroup.POST("/list", GoodsCategoryList)
		categoryGroup.POST("/del", GoodsCategoryDelete)
	}
}

func GoodsCategorySave(c *gin.Context) {
	var err error
	goodsCategoryLogic := logic.NewGoodsCategoryLogic(c)
	if err = c.ShouldBind(&goodsCategoryLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodsCategoryLogic.UID == "" {
		err = goodsCategoryLogic.Create()
		common.Response(c, err, goodsCategoryLogic)
		return
	}

	if err := goodsCategoryLogic.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodsCategoryLogic)
}

func GoodsCategoryList(c *gin.Context) {
	var body request.GoodsCategoryListReq
	if err := c.ShouldBind(&body); err != nil {
		common.Response(c, err, nil)
		return
	}
	var gclist []*logic.GoodsCategoryLogic
	var err error
	gclist, err = logic.NewGoodsCategoryLogic(c).ListGoodsCategoryByUser(body.Brief, body.LimitCond)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, gclist)
}

func GoodsCategoryDelete(c *gin.Context) {
	var err error
	goodCategoryLogic := logic.NewGoodsCategoryLogic(c)
	if err = c.ShouldBind(&goodCategoryLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodCategoryLogic.UID == "" {
		common.Response(c, common.RequestUIDMustErr, nil)
		return
	}

	err = goodCategoryLogic.Delete()
	common.Response(c, err, nil)
}

func GoodsList(c *gin.Context) {
	var err error
	var form request.GoodsListReq
	if err = c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	var goods []model.Goods
	goods, err = logic.NewGoodsLogic(c).LoadGoods(form)
	common.Response(c, err, goods)
}

func GoodsSave(c *gin.Context) {
	var err error
	goodLogic := logic.NewGoodsLogic(c)
	if err := c.ShouldBind(&goodLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodLogic.UID == "" {
		err = goodLogic.Create()
		common.Response(c, err, goodLogic)
		return
	}

	err = goodLogic.Update()
	common.Response(c, err, goodLogic)
}
