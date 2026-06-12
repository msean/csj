package logic

import (
	"app/global"
	"app/pkg/sms"
	"app/service/cache"
	"app/service/model/request"
	"fmt"

	"gorm.io/gorm"
)

func SmsKey(phone string) string {
	return fmt.Sprintf("csj:sms_code:%s", phone)
}

func SmsVerifyCodeSet(db *gorm.DB, phone string) (code string, err error) {
	return cache.SmsCache.SetVerifyCode(phone)
}

func SmsVerifyCodeCheck(db *gorm.DB, phone, input string) (right bool, err error) {
	return cache.SmsCache.VerifyCode(phone, input)
}

func SmsTodayCountCheck(db *gorm.DB, phone string) (over bool, err error) {
	return cache.SmsCache.CheckTodayCount(phone)
}

func SmsLoginAndRegister(sender sms.SmsSender, req request.SendVerifyCodeReq) (err error) {
	// over 判断是否超限
	var over bool
	if over, err = SmsTodayCountCheck(global.Global.DB, req.Phone); err != nil {
		return
	}
	if over {
		err = fmt.Errorf("发送验证码当日超限")
		return
	}
	var tempCode string
	// 注册
	if req.Typ == 1 {
		tempCode = global.Global.SmsRegisterTemp()
	}
	// 登陆
	if req.Typ == 2 {
		tempCode = global.Global.SmsLoginTemp()
	}

	// 设置验证码并存储
	var code string
	if code, err = SmsVerifyCodeSet(global.Global.DB, req.Phone); err != nil {
		return
	}

	msg := sms.SmsMsg{
		TemplateCode: tempCode,
		Mobile:       req.Phone,
		TemplateJson: map[string]any{
			"code": code,
		},
	}

	return sender.Send(msg)
}
