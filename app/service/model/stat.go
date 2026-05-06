package model

import (
	"app/global"
	"app/pkg/utils"
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

func BillingCondByOwnerUser(db *gorm.DB, owneruser string, customers []Customer) (billingLatestDate map[string]int, err error) {
	billingLatestDate = make(map[string]int)
	customUUIDList := make([]string, 0)
	for _, customer := range customers {
		billingLatestDate[customer.UID] = utils.DurationDays(customer.CreatedAt)
		customUUIDList = append(customUUIDList, customer.UID)
	}
	var results []struct {
		UserUUID      string
		LatestOrderAt time.Time
	}

	if err = db.Model(&BatchOrderGoods{}).Where("user_uuid in (?)", customUUIDList).
		Select("user_uuid, MAX(created_at) as latest_order_at").
		Group("user_uuid").
		Scan(&results).Error; err != nil {
		global.Global.Logger.Error(fmt.Sprintf("[BillingCondByOwnerUser] %s %s", owneruser, err))
		return
	}

	for _, result := range results {
		billingLatestDate[result.UserUUID] = utils.DurationDays(result.LatestOrderAt)
	}

	return
}
