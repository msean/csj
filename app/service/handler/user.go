package handler

import (
	"app/global"
	"app/service/common"
	"app/service/handler/middleware"
	"app/service/logic"
	"app/service/model"

	"github.com/gin-gonic/gin"
)

func userRouter(g *gin.RouterGroup) {
	userRouterGroup := g.Group("/user")
	{
		userRouterGroup.POST("/register", Register)
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
		"customerDebt":  creditAmount,
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
