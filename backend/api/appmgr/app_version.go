package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
	"go.uber.org/zap"
)

type AppVersionApi struct{}

// CreateAppVersion 创建App版本
// @Tags AppVersion
// @Summary 创建App版本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.AppVersion true "创建App版本"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /appVersion/createAppVersion [post]
func (api *AppVersionApi) CreateAppVersion(c *gin.Context) {
	var appVersion appmgr.AppVersion
	err := c.ShouldBindJSON(&appVersion)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = appVersionService.CreateAppVersion(&appVersion)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteAppVersion 删除App版本
// @Tags AppVersion
// @Summary 删除App版本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /appVersion/deleteAppVersion [delete]
func (api *AppVersionApi) DeleteAppVersion(c *gin.Context) {
	uid := c.Query("uid")
	err := appVersionService.DeleteAppVersion(uid)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAppVersionByIds 批量删除App版本
// @Tags AppVersion
// @Summary 批量删除App版本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /appVersion/deleteAppVersionByIds [delete]
func (api *AppVersionApi) DeleteAppVersionByIds(c *gin.Context) {
	uids := c.QueryArray("uids[]")
	err := appVersionService.DeleteAppVersionByIds(uids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAppVersion 更新App版本
// @Tags AppVersion
// @Summary 更新App版本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.AppVersion true "更新App版本"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /appVersion/updateAppVersion [put]
func (api *AppVersionApi) UpdateAppVersion(c *gin.Context) {
	var appVersion appmgr.AppVersion
	err := c.ShouldBindJSON(&appVersion)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = appVersionService.UpdateAppVersion(appVersion)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAppVersion 用uid查询App版本
// @Tags AppVersion
// @Summary 用uid查询App版本
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appmgr.AppVersion,msg=string} "查询成功"
// @Router /appVersion/findAppVersion [get]
func (api *AppVersionApi) FindAppVersion(c *gin.Context) {
	uid := c.Query("uid")
	res, err := appVersionService.GetAppVersion(uid)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// GetAppVersionList 分页获取App版本列表
// @Tags AppVersion
// @Summary 分页获取App版本列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appmgrReq.AppVersionSearch true "分页获取App版本列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /appVersion/getAppVersionList [get]
func (api *AppVersionApi) GetAppVersionList(c *gin.Context) {
	var pageInfo appmgrReq.AppVersionSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := appVersionService.GetAppVersionInfoList(pageInfo)
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
