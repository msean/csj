package handler

import (
	"app/global"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"
	"app/service/model/request"
	"app/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func userRouter(g *gin.RouterGroup) {
	userRouterGroup := g.Group("/user")
	{
		userRouterGroup.POST("/register", Register)
		userRouterGroup.POST("/send_verify_code", SendVerifyCode)
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

	amount, creditAmount, _ := model.MonthFinance(global.Global.DB, common.GetUserUUID(c))
	common.Response(c, nil, map[string]any{
		"name":          u.Name,
		"phone":         u.Phone,
		"customerDebt":  common.Float32Preserve(creditAmount, 2),
		"monthSales":    common.Float32Preserve(amount, 2),
		"vipRemainDays": 0,
	})
}

func Register(c *gin.Context) {

	var params request.RegisterParam
	if err := c.ShouldBind(&params); err != nil {
		common.Response(c, err, nil)
		return
	}

	userlogic := logic.NewUser(c)
	registerUser, e := userlogic.Register(params)
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	token, e := middleware.SetToken(params.Phone, registerUser.UID)
	if e != nil {
		common.Response(c, e, nil)
		return
	}

	common.Response(c, e, map[string]any{
		"token": token,
		"uuid":  utils.Violent2String(registerUser.UID),
	})
}

func SendVerifyCode(c *gin.Context) {
	var payload request.VerifyCodeParam
	var err error
	if err = c.ShouldBind(&payload); err != nil {
		common.Response(c, err, nil)
		return
	}
	// over 判断是否超限
	var over bool
	over, err = logic.SmsTodayCountCheck(global.Global.DB, payload.Phone)
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
		tempCode = global.Global.SmsRegisterTemp()
	}
	// 登陆
	if payload.Typ == 2 {
		tempCode = global.Global.SmsLoginTemp()
	}

	// 设置验证码并存储
	var code string
	if code, err = logic.SmsVerifyCodeSet(global.Global.DB, payload.Phone); err != nil {
		common.Response(c, err, nil)
		return
	}
	// 发送验证码
	if err = logic.SmsLoginAndRegister(global.Global.Sms, payload.Phone, code, tempCode); err != nil {
		common.Response(c, err, nil)
		return
	}
	common.Response(c, nil, nil)
}

func Login(c *gin.Context) {
	var param request.LoginParam
	if err := c.ShouldBind(&param); err != nil {
		common.Response(c, err, nil)
		return
	}
	userlogic := logic.NewUser(c)
	loginUser, e := userlogic.Login(param)
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	token, e := middleware.SetToken(loginUser.Phone, loginUser.UID)
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	common.Response(c, e, map[string]any{
		"token": token,
		"uuid":  strconv.FormatInt(loginUser.UID, 10),
	})
}

func UserUpdate(c *gin.Context) {
	var update request.UserUpdateParam
	if err := c.ShouldBind(&update); err != nil {
		common.Response(c, err, nil)
		return
	}

	userlogic := logic.NewUser(c)
	update.UID = common.GetUserUUID(c)
	e := userlogic.Update(update)
	if e != nil {
		common.Response(c, e, nil)
		return
	}
	common.Response(c, nil, nil)
}
