package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLogic struct {
	runTime *global.RunTime
	context *gin.Context
}

func NewUser(context *gin.Context) UserLogic {
	return UserLogic{
		runTime: global.Global,
		context: context,
	}
}

func (logic *UserLogic) CheckVerifyCode(phone, verifyCode string) (right bool, err error) {
	if verifyCode == "" {
		return false, nil
	}
	if logic.runTime.Env() == "test" || logic.runTime.Env() == "debug" {
		return verifyCode == logic.runTime.VerifyCode(), nil
	} else {
		right, err = SmsVerifyCodeCheck(logic.runTime.DB, phone, verifyCode)
		return
	}
}

func (logic *UserLogic) Register(params request.RegisterParam) (registerUser model.User, err error) {
	if params.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(params.Phone, params.VerifyCode); err != nil {
		return
	}
	if !right {
		err = common.VerifyCodeErr
		return
	}

	var user model.User
	err = model.Find(logic.runTime.DB, &user, model.NewWhereCond("phone", params.Phone))
	if err != nil {
		return
	}
	if user.UID != "" {
		// logic.runTime.Logger.Error(fmt.Sprintf("[UserLogic] [Register] phone: %s uid: %s", user.Phone, user.UID))
		logic.runTime.Logger.Error("[UserLogic] [Register]",
			zap.String("phone", user.Phone),
			zap.String("uid", user.UID))
		err = common.PhoneObejectExistErr
		return
	}
	registerUser = model.User{
		Phone: params.Phone,
	}
	tx := logic.runTime.DB.Begin()
	if err = model.CreateObj(tx, &registerUser); err != nil {
		logic.runTime.Logger.Error("[UserLogic] [Register] [CreateObj]",
			zap.String("phone", user.Phone),
			zap.Error(err))
		tx.Rollback()
		return
	}
	if err = model.NewTempCustomer(registerUser.UID, logic.runTime.DB); err != nil {
		logic.runTime.Logger.Error("[UserLogic] [Register] [NewTempCustomer]",
			zap.String("phone", registerUser.Phone),
			zap.String("UUID", registerUser.UID),
			zap.Error(err))
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func (logic *UserLogic) Login(params request.LoginParam) (loginUser model.User, err error) {
	if params.Phone == "" {
		err = common.PhoneCannotBeBlankErr
		return
	}
	var right bool
	if right, err = logic.CheckVerifyCode(params.Phone, params.VerifyCode); err != nil {
		return
	}
	if !right {
		err = common.VerifyCodeErr
		return
	}

	err = model.Find(logic.runTime.DB, &loginUser, model.NewWhereCond("phone", params.Phone))
	if err != nil {
		logic.runTime.Logger.Error("[UserLogic] [Login]",
			zap.String("phone", params.Phone),
			zap.Error(err))
		return
	}
	if loginUser.UID == "" {
		logic.runTime.Logger.Error("[UserLogic] [Login]",
			zap.String("phone", params.Phone))
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
		logic.runTime.Logger.Error("[UserLogic] [FromUUID]",
			zap.String("uuid", userUUID))
		err = common.UnRegisterErr
	}
	return
}

func (logic *UserLogic) Update(params request.UserUpdateParam) (err error) {
	var _u model.User
	if params.Phone != "" {
		e := model.Find(logic.runTime.DB, &_u, model.NewWhereCond("phone", params.Phone))
		if e != nil {
			err = e
			return
		}
		if _u.UID != "" && _u.UID != params.UID {
			err = common.PhoneObejectExistErr
			return
		}
	}
	return _u.Update(logic.runTime.DB)
}
