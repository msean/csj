package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
	"go.uber.org/zap"
)

type OfficialMessageApi struct{}

// CreateOfficialMessage 创建官方消息
// @Tags OfficialMessage
// @Summary 创建官方消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.OfficialMessage true "创建官方消息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /officialMessage/createOfficialMessage [post]
func (api *OfficialMessageApi) CreateOfficialMessage(c *gin.Context) {
	var officialMessage appmgr.OfficialMessage
	err := bindAppRecord(c, &officialMessage)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = officialMessageService.CreateOfficialMessage(&officialMessage)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteOfficialMessage 删除官方消息
// @Tags OfficialMessage
// @Summary 删除官方消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /officialMessage/deleteOfficialMessage [delete]
func (api *OfficialMessageApi) DeleteOfficialMessage(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.FailWithMessage("uid不能为空", c)
		return
	}
	err := officialMessageService.DeleteOfficialMessage(uid)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteOfficialMessageByIds 批量删除官方消息
// @Tags OfficialMessage
// @Summary 批量删除官方消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /officialMessage/deleteOfficialMessageByIds [delete]
func (api *OfficialMessageApi) DeleteOfficialMessageByIds(c *gin.Context) {
	uids := c.QueryArray("uids[]")
	err := officialMessageService.DeleteOfficialMessageByIds(uids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateOfficialMessage 更新官方消息
// @Tags OfficialMessage
// @Summary 更新官方消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.OfficialMessage true "更新官方消息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /officialMessage/updateOfficialMessage [put]
func (api *OfficialMessageApi) UpdateOfficialMessage(c *gin.Context) {
	var officialMessage appmgr.OfficialMessage
	err := bindAppRecord(c, &officialMessage)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = officialMessageService.UpdateOfficialMessage(officialMessage)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindOfficialMessage 用uid查询官方消息
// @Tags OfficialMessage
// @Summary 用uid查询官方消息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appmgr.OfficialMessage,msg=string} "查询成功"
// @Router /officialMessage/findOfficialMessage [get]
func (api *OfficialMessageApi) FindOfficialMessage(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		response.FailWithMessage("uid不能为空", c)
		return
	}
	res, err := officialMessageService.GetOfficialMessage(uid)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(presentAppRecord(res), c)
}

// GetOfficialMessageList 分页获取官方消息列表
// @Tags OfficialMessage
// @Summary 分页获取官方消息列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appmgrReq.OfficialMessageSearch true "分页获取官方消息列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /officialMessage/getOfficialMessageList [get]
func (api *OfficialMessageApi) GetOfficialMessageList(c *gin.Context) {
	var pageInfo appmgrReq.OfficialMessageSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if pageInfo.Page < 1 {
		pageInfo.Page = 1
	}
	if pageInfo.PageSize <= 0 {
		pageInfo.PageSize = 10
	}
	if pageInfo.PageSize > 100 {
		pageInfo.PageSize = 100
	}
	list, total, err := officialMessageService.GetOfficialMessageInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     presentAppRecord(list),
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
