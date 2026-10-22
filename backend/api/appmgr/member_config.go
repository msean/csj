package appmgr

import (
	"github.com/gin-gonic/gin"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"github.com/msean/csj/backend/model/common/response"
	"go.uber.org/zap"
)

type MemberConfigApi struct{}

// CreateMemberConfig 创建会员收费配置
// @Tags MemberConfig
// @Summary 创建会员收费配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.MemberConfig true "创建会员收费配置"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /memberConfig/createMemberConfig [post]
func (api *MemberConfigApi) CreateMemberConfig(c *gin.Context) {
	var memberConfig appmgr.MemberConfig
	err := c.ShouldBindJSON(&memberConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = memberConfigService.CreateMemberConfig(&memberConfig)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteMemberConfig 删除会员收费配置
// @Tags MemberConfig
// @Summary 删除会员收费配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /memberConfig/deleteMemberConfig [delete]
func (api *MemberConfigApi) DeleteMemberConfig(c *gin.Context) {
	uid := c.Query("uid")
	err := memberConfigService.DeleteMemberConfig(uid)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteMemberConfigByIds 批量删除会员收费配置
// @Tags MemberConfig
// @Summary 批量删除会员收费配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /memberConfig/deleteMemberConfigByIds [delete]
func (api *MemberConfigApi) DeleteMemberConfigByIds(c *gin.Context) {
	uids := c.QueryArray("uids[]")
	err := memberConfigService.DeleteMemberConfigByIds(uids)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateMemberConfig 更新会员收费配置
// @Tags MemberConfig
// @Summary 更新会员收费配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body appmgr.MemberConfig true "更新会员收费配置"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /memberConfig/updateMemberConfig [put]
func (api *MemberConfigApi) UpdateMemberConfig(c *gin.Context) {
	var memberConfig appmgr.MemberConfig
	err := c.ShouldBindJSON(&memberConfig)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = memberConfigService.UpdateMemberConfig(memberConfig)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindMemberConfig 用uid查询会员收费配置
// @Tags MemberConfig
// @Summary 用uid查询会员收费配置
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=appmgr.MemberConfig,msg=string} "查询成功"
// @Router /memberConfig/findMemberConfig [get]
func (api *MemberConfigApi) FindMemberConfig(c *gin.Context) {
	uid := c.Query("uid")
	res, err := memberConfigService.GetMemberConfig(uid)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

// GetMemberConfigList 分页获取会员收费配置列表
// @Tags MemberConfig
// @Summary 分页获取会员收费配置列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query appmgrReq.MemberConfigSearch true "分页获取会员收费配置列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /memberConfig/getMemberConfigList [get]
func (api *MemberConfigApi) GetMemberConfigList(c *gin.Context) {
	var pageInfo appmgrReq.MemberConfigSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := memberConfigService.GetMemberConfigInfoList(pageInfo)
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
