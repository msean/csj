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
	var param request.GoodsCategorySaveParam
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodsCategoryLogic := logic.NewGoodsCategoryLogic(c)
	var goodsCategory model.GoodsCategory
	var err error
	if param.UID == "" {
		if err := goodsCategoryLogic.Check(param); err != nil {
			common.Response(c, err, nil)
			return
		}
		if goodsCategory, err = goodsCategoryLogic.Create(param); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, goodsCategoryLogic)
		return
	}

	if goodsCategory, err = goodsCategoryLogic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodsCategory)
}

func GoodsCategoryList(c *gin.Context) {
	var params request.ListGoodsCategoryParam
	if err := c.ShouldBind(&params); err != nil {
		common.Response(c, err, nil)
		return
	}
	var gclist []*model.GoodsCategory
	var err error
	if gclist, err = logic.NewGoodsCategoryLogic(c).ListGoodsCategoryByUser(params.Brief, params.LimitCond); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, gclist)
}

func GoodsCategoryDelete(c *gin.Context) {
	var param request.DeleteGoodsCategoryParam
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodCategoryLogic := logic.NewGoodsCategoryLogic(c)
	if param.UID == "" {
		common.Response(c, common.RequestUIDMustErr, nil)
		return
	}

	if err := goodCategoryLogic.Delete(param.UID); err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, nil)
}

func GoodsList(c *gin.Context) {
	var form request.ListGoodsParam
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	goods, err := logic.NewGoodsLogic(c).LoadGoods(form.SearchKey, form.LimitCond)

	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, goods)
}

func GoodsSave(c *gin.Context) {
	var param request.GoodsSaveParam

	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}

	goodLogic := logic.NewGoodsLogic(c)
	var goods model.Goods
	var err error
	if param.UID == "" {
		if err := goodLogic.Check(param); err != nil {
			common.Response(c, err, nil)
			return
		}
		if goods, err = goodLogic.Create(param); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, goodLogic)
		return
	}

	if goods, err = goodLogic.Update(param); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goods)
}
