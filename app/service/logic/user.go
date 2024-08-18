package logic

import (
	"app/global"
	"app/service/common"
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
		runTime: global.GlobalRunTime,
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

	err = model.Find(logic.runTime.DB, &logic.User, logic.WherePhoneCond())
	if err != nil {
		return
	}
	if logic.UID != "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Register] phone: %s uid: %s", logic.Phone, logic.UID))
		err = common.PhoneObejectExistErr
		return
	}
	user := model.User{
		Phone: logic.Phone,
	}
	tx := logic.runTime.DB.Begin()
	if err = model.CreateObj(tx, &user); err != nil {
		tx.Rollback()
		return
	}
	if err = model.NewTempCustomer(user.UID, logic.runTime.DB); err != nil {
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

	err = model.Find(logic.runTime.DB, &logic.User, logic.WherePhoneCond())
	if err != nil {
		return
	}
	if logic.UID == "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Login] phone: %s", logic.Phone))
		err = common.PhoneUnRegisterErr
		return
	}
	return
}

func (logic *UserLogic) FromUUID(userUUID string) (user model.User, err error) {
	err = model.Find(logic.runTime.DB, &user, model.WhereUIDCond(userUUID))
	if err != nil {
		return
	}
	if user.UID == "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Login] uid: %s", logic.UID))
		err = common.UnRegisterErr
	}
	return
}

func (logic *UserLogic) Update() (err error) {
	if logic.Phone != "" {
		var _u model.User
		e := model.Find(logic.runTime.DB, &_u, logic.WherePhoneCond())
		if e != nil {
			err = e
			return
		}
		if _u.UID != "" && _u.UID != logic.UID {
			err = common.PhoneObejectExistErr
			return
		}
	}
	return logic.User.Update(logic.runTime.DB)
}
