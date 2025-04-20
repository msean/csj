package biz

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type GoodsRoter struct{}

func (s *GoodsRoter) InitGoodsRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	router := Router.Group("goods").Use(middleware.OperationRecord())

	var api = v1.ApiGroupApp.BizApiGroup.GoodsApi
	{
		router.GET("list_goods", api.ListGoods) // 新建users表
	}
}
