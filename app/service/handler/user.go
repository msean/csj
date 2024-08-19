package handler

import (
	"app/global"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"fmt"

	"github.com/gin-gonic/gin"
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
	_u := logic.NewUser(c)
	u, err := _u.FromUUID(common.GetUserUUID(c))

	if err != nil {
		common.Response(c, err, nil)
		return
	}

	amount, creditAmount, _ := model.MonthFinance(global.GlobalRunTime.DB, common.GetUserUUID(c))
	common.Response(c, nil, map[string]any{
		"name":          u.Name,
		"phone":         u.Phone,
		"customerDebt":  common.Float32Preserve(creditAmount, 2),
		"monthSales":    common.Float32Preserve(amount, 2),
		"vipRemainDays": 0,
	})
}

func Register(c *gin.Context) {
	userlogic := logic.NewUser(c)
	if err := c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	e := userlogic.Register()
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	token, e := middleware.SetToken(userlogic.Phone, userlogic.UID)
	if e != nil {
		common.Response(c, e, nil)
		return
	}

	common.Response(c, e, map[string]any{
		"token": token,
		"uuid":  userlogic.UID,
	})
}

func SenderVerifyCode(c *gin.Context) {
	type Payload struct {
		Phone string `json:"phone"`
		Typ   int    `json:"type"`
	}
	var payload Payload
	var err error
	if err = c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}
	// over 判断是否超限
	var over bool
	over, err = logic.SmsTodayCountCheck(global.GlobalRunTime.DB, payload.Phone)
	if err != nil {
		common.Response(c, err, nil)
		return
	}
	if over {
		common.Response(c, fmt.Errorf("发送验证码当日超限"), nil)
		return
	}
	var tempCode string
	// 注册
	if payload.Typ == 1 {
		tempCode = global.GlobalRunTime.SmsRegisterTemp()
	}
	// 登陆
	if payload.Typ == 2 {
		tempCode = global.GlobalRunTime.SmsLoginTemp()
	}

	// 设置验证码并存储
	var code string
	if code, err = logic.SmsVerifyCodeSet(global.GlobalRunTime.DB, payload.Phone); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 发送验证码
	if err = logic.SmsLoginAndRegister(global.GlobalRunTime.Sms, payload.Phone, code, tempCode); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func Login(c *gin.Context) {
	userlogic := logic.NewUser(c)
	if err := c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	e := userlogic.Login()
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	token, e := middleware.SetToken(userlogic.Phone, userlogic.UID)
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	common.Response(c, e, map[string]any{
		"token": token,
		"uuid":  userlogic.UID,
	})
}

func UserUpdate(c *gin.Context) {
	userlogic := logic.NewUser(c)
	if err := c.ShouldBind(&userlogic); err != nil {
		common.Response(c, err, nil)
		return
	}
	userlogic.UID = common.GetUserUUID(c)
	e := userlogic.Update()
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	common.Response(c, nil, userlogic.User)
}
