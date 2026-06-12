package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
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

func (logic *UserLogic) Register() (err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(); err != nil {
		return
	}
	if !right {
		return common.VerifyCodeErr
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
	return
}

func (logic *UserLogic) Login() (err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(); err != nil {
		return
	}
	if !right {
		return common.VerifyCodeErr
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
