package dao

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type orderDao struct{}

func newOrderDao() *orderDao {
	return &orderDao{}
}

func (dao *orderDao) UpdateStatus(db *gorm.DB, orderUUID string, status int) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Model(&model.BatchOrder{}).Update("status", status).Error
}

func (dao *orderDao) Shared(db *gorm.DB, orderUUID string) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Updates(&model.BatchOrder{
		Shared: common.BatchOrderShared,
	}).Error
}

func (dao *orderDao) UpdateGoods(db *gorm.DB, object model.BatchOrderGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderGoods{
		Price:  object.Price,
		Mount:  object.Mount,
		Weight: object.Weight,
	}).Error
}

func (dao *orderDao) UpdateOrderPay(db *gorm.DB, object model.BatchOrderPay) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderPay{
		PayType: object.PayType,
		Amount:  object.Amount,
	}).Error
}

func (dao *orderDao) LatestOrderByCustomers(db *gorm.DB, owneruser string, customers []model.Customer) (billingLatestDate map[string]int, err error) {
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

	if err = db.Model(&model.BatchOrderGoods{}).Where("user_uuid in (?)", customUUIDList).
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

func (dao *orderDao) MonthFinance(db *gorm.DB, ownerUser string) (amount, creditAmount float64, err error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var bos []model.BatchOrder
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

func (dao *orderDao) ListByBatchUUIDIn(db *gorm.DB, batchUUID string, goodsUUIDList []string) (objects []model.BatchOrderGoods, err error) {
	err = utils.Find(db, objects, utils.NewWhereCond("batch_uuid", batchUUID), utils.NewInCondFromString("goods_uuid", goodsUUIDList))
	return
}
