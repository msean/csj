package biz

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type UsersRouter struct{}

// InitUsersRouter 初始化 users表 路由信息
func (s *UsersRouter) InitUsersRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	usersRouter := Router.Group("users").Use(middleware.OperationRecord())
	usersRouterWithoutRecord := Router.Group("users")
	usersRouterWithoutAuth := PublicRouter.Group("users")

	var usersApi = v1.ApiGroupApp.BizApiGroup.UsersApi
	{
		usersRouter.POST("createUsers", usersApi.CreateUsers)             // 新建users表
		usersRouter.DELETE("deleteUsers", usersApi.DeleteUsers)           // 删除users表
		usersRouter.DELETE("deleteUsersByIds", usersApi.DeleteUsersByIds) // 批量删除users表
		usersRouter.PUT("updateUsers", usersApi.UpdateUsers)              // 更新users表
	}
	{
		usersRouterWithoutRecord.GET("findUsers", usersApi.FindUsers)       // 根据ID获取users表
		usersRouterWithoutRecord.GET("getUsersList", usersApi.GetUsersList) // 获取users表列表
	}
	{
		usersRouterWithoutAuth.GET("getUsersPublic", usersApi.GetUsersPublic) // 获取users表列表
	}
}
