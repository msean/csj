package dao

import (
	"app/pkg/utils"
	"app/service/model"
	"time"

	"gorm.io/gorm"
)

// PaymentDao 支付相关DAO
type paymentDao struct{}

func newPaymentDao() *paymentDao {
	return &paymentDao{}
}

// CreditSummary 赊欠汇总结构
type CreditSummary struct {
	UserUUID    string  `gorm:"column:user_uuid"`
	TotalCredit float64 `gorm:"column:total_credit"`
}

// CreatePay 创建还款记录（分表）
func (dao *paymentDao) CreatePay(db *gorm.DB, pay *model.BatchOrderPay) error {
	return db.Table(ShardPayTable(pay.OwnerUser)).Create(pay).Error
}

// FindPayRecords 查询还款记录
func (dao *paymentDao) FindPayRecords(db *gorm.DB, ownerUser, customerUUID string, offset, limit int) ([]model.BatchOrderPay, int64, error) {
	var pays []model.BatchOrderPay
	var total int64

	query := db.Table(ShardPayTable(ownerUser)).
		Where("owner_user = ? AND customer_uuid = ? and is_revoked=0", ownerUser, customerUUID).
		Order("created_at DESC")

	// 总数
	query.Count(&total)

	// 分页
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Find(&pays).Error
	return pays, total, err
}

// FindPayByUID 查询单条还款记录
func (dao *paymentDao) FindPayByUID(db *gorm.DB, ownerUser, payUUID string) (pay *model.BatchOrderPay, err error) {
	err = db.Table(ShardPayTable(ownerUser)).Where("uid = ? AND owner_user = ?", payUUID, ownerUser).First(&pay).Error
	if err != nil {
		return nil, err
	}
	return
}

// CreateMessage 创建消息记录
func (dao *paymentDao) CreateMessage(db *gorm.DB, message *model.MessageCenter) error {
	return db.Create(message).Error
}

// FindMessages 查询消息列表
func (dao *paymentDao) FindMessages(db *gorm.DB, ownerUser string, offset, limit int) ([]model.MessageCenter, error) {
	var messages []model.MessageCenter
	query := db.Model(&model.MessageCenter{}).
		Where("owner_user = ?", ownerUser).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Find(&messages).Error
	return messages, err
}

// CountMessages 统计消息总数和未读数
func (dao *paymentDao) CountMessages(db *gorm.DB, ownerUser string) (totalCount, unreadCount int64, err error) {
	err = db.Model(&model.MessageCenter{}).Where("owner_user = ?", ownerUser).Count(&totalCount).Error
	if err != nil {
		return
	}
	err = db.Model(&model.MessageCenter{}).Where("owner_user = ? AND is_read = 0", ownerUser).Count(&unreadCount).Error
	return
}

// MarkMessageRead 标记消息为已读
func (dao *paymentDao) MarkMessageRead(db *gorm.DB, ownerUser, messageUUID string) error {
	return db.Model(&model.MessageCenter{}).
		Where("uid = ? AND owner_user = ?", messageUUID, ownerUser).
		Update("is_read", 1).Error
}

// UpdateOrderPay 更新订单支付（分表）
func (dao *paymentDao) UpdateOrderPay(db *gorm.DB, ownerUser string, object model.BatchOrderPay) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Table(ShardPayTable(ownerUser)).Updates(&model.BatchOrderPay{
		PayType: object.PayType,
		Amount:  object.Amount,
	}).Error
}

// UpdatePayRevokeStatus 更新还款记录的撤销状态（分表）
func (dao *paymentDao) UpdatePayRevokeStatus(db *gorm.DB, ownerUser, payUUID string, revokedAt time.Time, reason string) error {
	return db.Table(ShardPayTable(ownerUser)).Where("uid = ?", payUUID).
		Updates(map[string]interface{}{
			"is_revoked":     1,
			"revoked_at":     revokedAt,
			"revoked_reason": reason,
		}).Error
}
