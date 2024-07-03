package handler

import (
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"

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
	goodsCategoryLogic := logic.NewGoodsCategoryLogic(c)
	if err := c.ShouldBind(&goodsCategoryLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodsCategoryLogic.UID == "" {
		if err := goodsCategoryLogic.Check(); err != nil {
			common.Response(c, err, nil)
			return
		}
		if err := goodsCategoryLogic.Create(); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, goodsCategoryLogic)
		return
	}

	if err := goodsCategoryLogic.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodsCategoryLogic)
}

func GoodsCategoryList(c *gin.Context) {
	type Body struct {
		Brief bool `json:"brief"`
		model.LimitCond
	}
	var body Body
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

	goodCategoryLogic := logic.NewGoodsCategoryLogic(c)
	if err := c.ShouldBind(&goodCategoryLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodCategoryLogic.UID == "" {
		common.Response(c, common.RequestUIDMustErr, nil)
		return
	}

	if err := goodCategoryLogic.Delete(); err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, nil)
}

func GoodsList(c *gin.Context) {

	type Form struct {
		SearchKey string `json:"searchName"`
		model.LimitCond
	}
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		common.Response(c, err, nil)
		return
	}

	goods, err := logic.NewGoodsLogic(c).LoadGoods(common.GetUserUUID(c), form.SearchKey, form.LimitCond)

	if err != nil {
		common.Response(c, err, nil)
		return
	}

	common.Response(c, nil, goods)
}

func GoodsSave(c *gin.Context) {

	goodLogic := logic.NewGoodsLogic(c)
	if err := c.ShouldBind(&goodLogic); err != nil {
		common.Response(c, err, nil)
		return
	}

	if goodLogic.UID == "" {
		if err := goodLogic.Check(); err != nil {
			common.Response(c, err, nil)
			return
		}
		if err := goodLogic.Create(); err != nil {
			common.Response(c, err, nil)
			return
		}
		common.Response(c, nil, goodLogic)
		return
	}

	if err := goodLogic.Update(); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, goodLogic)
}
