package model

import (
	"app/global"
	"app/service/common"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func MonthFinance(db *gorm.DB, ownerUser string) (amount, creditAmount float64, err error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var bos []BatchOrder
	if err = db.
		Where("owner_user=?", ownerUser).
		Where("created_at>=?", monthStart).
		Where("status not in (?)", common.FinalBatchOrder).
		Find(&bos).Error; err != nil {
		global.Global.Logger.Error(fmt.Sprintf("[MonthFinance] %s %s", ownerUser, err))
	}
	for _, bo := range bos {
		amount += bo.TotalAmount
		creditAmount += bo.CreditAmount
	}
	return
}
