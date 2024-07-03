package model

import (
	"app/service/common"
	"time"

	"gorm.io/gorm"
)

func MonthAmountByOwnerUser(db *gorm.DB, ownerUser string) (amount float32, err error) {
	type OrderGoods struct {
		Price  float32
		Weight float32
		Mount  int32
	}

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var orderList []string
	if err = db.Model(&BatchOrder{}).
		Where("owner_user=?", ownerUser).
		Where("created_at>=?", monthStart).
		Where("status not in (?)", []int32{BatchOrderReTurn, BatchOrderCancel, BatchOrderRefund}).
		Select("uid").
		Pluck("uid", &orderList).Error; err != nil {
		return
	}
	var orderGoods []BatchGoods
	if err = db.Model(&BatchOrderGoods{}).Where("batch_order_uuid in (?)", orderList).Scan(&orderGoods).Error; err != nil {
		return
	}
	if len(orderGoods) == 0 {
		return
	}

	for _, order := range orderGoods {
		if order.Weight == 0 {
			amount += order.Price * float32(order.Mount)
		} else {
			amount += order.Price * order.Weight
		}
	}

	return
}

func BillingCondByOwnerUser(db *gorm.DB, owneruser string, customers []Customer) (billingLatestDate map[string]int, err error) {
	billingLatestDate = make(map[string]int)
	customUUIDList := make([]string, 0)
	for _, customer := range customers {
		billingLatestDate[customer.UID] = common.DurationDays(customer.CreatedAt)
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
		return
	}

	for _, result := range results {
		billingLatestDate[result.UserUUID] = common.DurationDays(result.LatestOrderAt)
	}

	return
}
