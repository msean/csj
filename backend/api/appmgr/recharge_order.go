package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
	"go.uber.org/zap"
)

type RechargeOrderApi struct{}

// FindRechargeOrder 用uid查询充值订单
// @Tags RechargeOrder
// @Summary 用uid查询充值订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appmgr.RechargeOrder,msg=string} "查询成功"
// @Router /rechargeOrder/findRechargeOrder [get]
func (api *RechargeOrderApi) FindRechargeOrder(c *gin.Context) {
	uid := c.Query("uid")
	res, err := rechargeOrderService.GetRechargeOrder(uid)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// GetRechargeOrderList 分页获取充值订单列表
// @Tags RechargeOrder
// @Summary 分页获取充值订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appmgrReq.RechargeOrderSearch true "分页获取充值订单列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /rechargeOrder/getRechargeOrderList [get]
func (api *RechargeOrderApi) GetRechargeOrderList(c *gin.Context) {
	var pageInfo appmgrReq.RechargeOrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := rechargeOrderService.GetRechargeOrderInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
