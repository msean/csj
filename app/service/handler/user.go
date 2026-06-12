package handler

import (
	"app/global"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func userRouter(g *gin.RouterGroup) {
	userRouterGroup := g.Group("/user")
	{
		userRouterGroup.POST("/register", Register)
		userRouterGroup.POST("/send_verify_code", SenderVerifyCode)
		userRouterGroup.POST("/login", Login)
		userRouterGroup.POST("/update", middleware.AuthMiddleware(), UserUpdate)
		userRouterGroup.GET("/profile", middleware.AuthMiddleware(), UserInfo)
	}
}

func UserInfo(c *gin.Context) {
	logic := logic.NewUser(c)
	rsp, err := logic.Profile(common.GetUserUUID(c))
	common.Response(c, err, rsp)
}

func Register(c *gin.Context) {
	userlogic := logic.NewUser(c)
	if err := c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	var err error
	var token string
	token, err = userlogic.Register()
	common.Response(c, err, map[string]any{
		"token": token,
		"uuid":  userlogic.UID,
	})
}

func SenderVerifyCode(c *gin.Context) {
	var payload request.SendVerifyCodeReq
	var err error
	if err = c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 发送验证码
	if err = logic.SmsLoginAndRegister(global.Global.Sms, payload); err != nil {
		global.Global.Logger.Error("SenderVerifyCode SmsLoginAndRegister", zap.Any("request", payload), zap.Error(err))
	}
	common.Response(c, err, nil)
}

func Login(c *gin.Context) {
	userlogic := logic.NewUser(c)
	if err := c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	var err error
	var token string
	token, err = userlogic.Login()

	common.Response(c, err, map[string]any{
		"token": token,
		"uuid":  userlogic.UID,
	})
}

func UserUpdate(c *gin.Context) {
	var err error
	userlogic := logic.NewUser(c)
	if err = c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	err = userlogic.Update()
	common.Response(c, err, userlogic.User)
}
