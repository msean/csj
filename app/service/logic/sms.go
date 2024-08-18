package logic

import (
	"app/pkg"
	"app/service/model"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

func SmsKey(phone string) string {
	return fmt.Sprintf("csj:sms_code:%s", phone)
}

// func SmsVerifyCodeSet(phone string) (code string, err error) {
// 	rand.Seed(time.Now().UnixNano())
// 	code = fmt.Sprintf("%d", rand.Intn(900000)+100000)
// 	err = global.GlobalRunTime.Redis.Set(context.Background(), SmsKey(phone), code, 60*time.Second).Err()
// 	return
// }

// func SmsVerifyCodeCheck(phone, input string) bool {
// 	return global.GlobalRunTime.Redis.Get(context.Background(), SmsKey(phone)).String() == input
// }

func SmsVerifyCodeSet(db *gorm.DB, phone string) (code string, err error) {
	rand.Seed(time.Now().UnixNano())
	code = fmt.Sprintf("%d", rand.Intn(900000)+100000)
	err = model.CreateObj(db, model.Sms{
		Phone: phone,
		Code:  code,
	})
	return
}

func SmsVerifyCodeCheck(db *gorm.DB, phone, input string) (right bool, err error) {
	var sms model.Sms
	if err = model.Find(db, &sms, model.NewWhereCond("phone", phone), model.NewOrderCond("created_at desc")); err != nil {
		return
	}
	right = (sms.Code == input)
	return
}

func SmsTodayCountCheck(db *gorm.DB, phone string) (over bool, err error) {
	var count int64
	todayStart := time.Now().Truncate(24 * time.Hour)
	if err = db.Model(&model.Sms{}).Where("phone=? andcreated_at >= ?", phone, todayStart).Count(&count).Error; err != nil {
		return
	}
	over = count <= 5
	return
}

func SmsLoginAndRegister(sender pkg.SmsSender, phone, code, templateCode string) {
	msg := pkg.SmsMsg{
		TemplateCode: templateCode,
		Mobile:       phone,
		TemplateJson: map[string]any{
			"code": code,
		},
	}
	sender.Send(msg)
}
