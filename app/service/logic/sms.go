package logic

import (
	"app/pkg/sms"
	"app/service/cache"
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

func SmsLoginAndRegister(sender sms.SmsSender, phone, code, templateCode string) error {
	msg := sms.SmsMsg{
		TemplateCode: templateCode,
		Mobile:       phone,
		TemplateJson: map[string]any{
			"code": code,
		},
	}

	return sender.Send(msg)
}
