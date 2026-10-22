package dao

import (
	"app/service/model"
	"time"

	"gorm.io/gorm"
)

type rechargeDao struct{}

func newRechargeDao() *rechargeDao {
	return &rechargeDao{}
}

// Create 创建充值订单
func (dao *rechargeDao) Create(db *gorm.DB, order *model.RechargeOrder) error {
	return db.Create(order).Error
}

// FindByOrderNo 根据订单号查询
func (dao *rechargeDao) FindByOrderNo(db *gorm.DB, orderNo string) (model.RechargeOrder, error) {
	var order model.RechargeOrder
	err := db.Where("order_no = ?", orderNo).First(&order).Error
	return order, err
}

// MarkPaid 标记订单已支付（幂等：仅更新未支付订单）
func (dao *rechargeDao) MarkPaid(db *gorm.DB, orderNo, transactionID string) (model.RechargeOrder, error) {
	var order model.RechargeOrder
	err := db.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return order, err
	}
	// 已支付直接返回，保证幂等
	if order.PayStatus == 1 {
		return order, nil
	}
	now := time.Now()
	updates := map[string]any{
		"pay_status":     1,
		"pay_time":       now,
		"transaction_id": transactionID,
	}
	err = db.Model(&order).Where("order_no = ? AND pay_status != ?", orderNo, 1).Updates(updates).Error
	if err != nil {
		return order, err
	}
	order.PayStatus = 1
	order.PayTime = &now
	order.TransactionID = transactionID
	return order, nil
}

// ListByUser 按用户分页查询充值记录
func (dao *rechargeDao) ListByUser(db *gorm.DB, userUUID string, page, pageCount int) ([]model.RechargeOrder, int64, error) {
	var list []model.RechargeOrder
	var total int64
	q := db.Model(&model.RechargeOrder{}).Where("user_uuid = ?", userUUID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if pageCount > 0 {
		q = q.Limit(pageCount).Offset((page - 1) * pageCount)
	}
	err := q.Order("created_at DESC").Find(&list).Error
	return list, total, err
}
