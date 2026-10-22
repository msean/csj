package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
	"go.uber.org/zap"
)

type AppUserApi struct{}

// FindAppUser 用uid查询用户
// @Tags AppUser
// @Summary 用uid查询用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appmgr.AppUser,msg=string} "查询成功"
// @Router /appUser/findAppUser [get]
func (api *AppUserApi) FindAppUser(c *gin.Context) {
	uid := c.Query("uid")
	res, err := appUserService.GetAppUser(uid)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// GetAppUserList 分页获取用户列表
// @Tags AppUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appmgrReq.AppUserSearch true "分页获取用户列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /appUser/getAppUserList [get]
func (api *AppUserApi) GetAppUserList(c *gin.Context) {
	var pageInfo appmgrReq.AppUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := appUserService.GetAppUserInfoList(pageInfo)
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
