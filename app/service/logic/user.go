package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/handler/middleware"
	"app/service/model"
	"app/service/model/response"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserLogic struct {
	runTime *global.RunTime
	context *gin.Context
	model.User
	VerfifyCode string `json:"verifycode"`
}

func NewUser(context *gin.Context) UserLogic {
	return UserLogic{
		runTime: global.Global,
		context: context,
	}
}

func (logic *UserLogic) CheckVerifyCode() (right bool, err error) {
	if logic.VerfifyCode == "" {
		return false, nil
	}
	if logic.runTime.Env() == "test" || logic.runTime.Env() == "debug" {
		return logic.VerfifyCode == logic.runTime.VerifyCode(), nil
	} else {
		right, err = SmsVerifyCodeCheck(logic.runTime.DB, logic.Phone, logic.VerfifyCode)
		return
	}
}

func (logic *UserLogic) Register() (token string, err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(); err != nil {
		return
	}
	if !right {
		err = common.VerifyCodeErr
		return
	}

	var userModel model.User
	if userModel, err = dao.UserDao.FindByPhone(logic.runTime.DB, logic.Phone); err != nil {
		return
	}
	if userModel.UID != "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Register] phone: %s uid: %s", logic.Phone, logic.UID))
		err = common.PhoneObejectExistErr
		return
	}
	user := model.User{
		Phone: logic.Phone,
	}
	tx := logic.runTime.DB.Begin()
	if err = utils.CreateObj(tx, &user); err != nil {
		tx.Rollback()
		return
	}
	if err = dao.CustomerDao.NewTempCustomer(tx, user.UID); err != nil {
		tx.Rollback()
		return
	}
	logic.User = user
	tx.Commit()

	token, err = middleware.SetToken(logic.User.Phone, logic.User.UID)
	return
}

func (logic *UserLogic) Login() (token string, err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(); err != nil {
		return
	}
	if !right {
		err = common.VerifyCodeErr
		return
	}

	var userModel model.User
	if userModel, err = dao.UserDao.FindByPhone(logic.runTime.DB, logic.Phone); err != nil {
		return
	}
	if userModel.UID == "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Login] phone: %s", logic.Phone))
		err = common.PhoneUnRegisterErr
		return
	}
	token, err = middleware.SetToken(logic.Phone, logic.UID)
	return
}

func (logic *UserLogic) Update() (err error) {
	if logic.Phone != "" {
		var _u model.User
		if _u, err = dao.UserDao.FindByPhone(logic.runTime.DB, logic.Phone); err != nil {
			return
		}
		if _u.UID != "" && _u.UID != logic.UID {
			err = common.PhoneObejectExistErr
			return
		}
	}
	return dao.UserDao.Update(logic.runTime.DB, logic.User)
}

func (logic *UserLogic) Profile(userUID string) (rsp response.UserProfileRsp, err error) {
	var user model.User
	if user, err = dao.UserDao.FromUUID(global.Global.DB, userUID); err != nil {
		return
	}

	amount, _ := dao.OrderDao.MonthSales(global.Global.DB, userUID)
	creditAmount, _ := dao.OrderDao.CreditAmountTotal(global.Global.DB, userUID)
	rsp.User = user
	rsp.CreditAmount = utils.FloatReserveStr(creditAmount, 1)
	rsp.MonthSales = utils.FloatReserveStr(amount, 1)
	return
}
