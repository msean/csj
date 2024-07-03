package user

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	userReq "github.com/flipped-aurora/gin-vue-admin/server/model/user/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UsersApi struct{}

var usersService = service.ServiceGroupApp.UserServiceGroup.UsersService

func (usersApi *UsersApi) CreateUsers(c *gin.Context) {
	var users user.Users
	err := c.ShouldBindJSON(&users)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := usersService.CreateUsers(&users); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func (usersApi *UsersApi) DeleteUsers(c *gin.Context) {
	ID := c.Query("ID")
	if err := usersService.DeleteUsers(ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func (usersApi *UsersApi) DeleteUsersByIds(c *gin.Context) {
	IDs := c.QueryArray("IDs[]")
	if err := usersService.DeleteUsersByIds(IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateUsers 更新users表
// @Tags Users
// @Summary 更新users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body user.Users true "更新users表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /users/updateUsers [put]
func (usersApi *UsersApi) UpdateUsers(c *gin.Context) {
	var users user.Users
	err := c.ShouldBindJSON(&users)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := usersService.UpdateUsers(users); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindUsers 用id查询users表
// @Tags Users
// @Summary 用id查询users表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query user.Users true "用id查询users表"
// @Success 200 {object} response.Response{data=object{reusers=user.Users},msg=string} "查询成功"
// @Router /users/findUsers [get]
func (usersApi *UsersApi) FindUsers(c *gin.Context) {
	ID := c.Query("UID")
	if reusers, err := usersService.GetUsers(ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reusers, c)
	}
}

// GetUsersList 分页获取users表列表
// @Tags Users
// @Summary 分页获取users表列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query userReq.UsersSearch true "分页获取users表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /users/getUsersList [get]
func (usersApi *UsersApi) GetUsersList(c *gin.Context) {
	var pageInfo userReq.UsersSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := usersService.GetUsersInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetUsersPublic 不需要鉴权的users表接口
// @Tags Users
// @Summary 不需要鉴权的users表接口
// @accept application/json
// @Produce application/json
// @Param data query userReq.UsersSearch true "分页获取users表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /users/getUsersPublic [get]
func (usersApi *UsersApi) GetUsersPublic(c *gin.Context) {
	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的users表接口信息",
	}, "获取成功", c)
}
