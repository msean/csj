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

func (logic *UserLogic) CheckVerifyCode() bool {
	if logic.runTime.Env() == "test" || logic.runTime.Env() == "debug" {
		return logic.VerfifyCode == logic.runTime.VerifyCode()
	}
	return false
}

func (logic *UserLogic) Register() (err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	if !logic.CheckVerifyCode() {
		err = common.VerifyCodeErr
		return
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
	err = model.CreateObj(logic.runTime.DB, &user)
	logic.User = user
	return
}

func (logic *UserLogic) Login() (err error) {
	if logic.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	if !logic.CheckVerifyCode() {
		err = common.VerifyCodeErr
		return
	}

	err = model.Find(logic.runTime.DB, &logic.User, logic.WherePhoneCond())
	if err != nil {
		return
	}
	if logic.UID == "" {
		logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Login] phone: %s", logic.Phone))
		err = common.PhoneUnRegisterErr
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
	err = logic.User.Update(logic.runTime.DB)
	return
}
