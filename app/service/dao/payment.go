package dao

import (
	"app/service/model"
	"fmt"

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

// FindCreditSummaries 查询所有有赊欠的客户汇总
func (dao *paymentDao) FindCreditSummaries(db *gorm.DB, ownerUser string, validStatuses []int) ([]CreditSummary, error) {
	var summaries []CreditSummary
	err := db.Model(&model.BatchOrder{}).
		Select("user_uuid, SUM(credit_amount) as total_credit").
		Where("owner_user = ? AND status IN (?) AND credit_amount != 0", ownerUser, validStatuses).
		Group("user_uuid").
		Scan(&summaries).Error
	return summaries, err
}

// FindCreditSummariesWithPage 分页查询有赊欠的客户汇总
func (dao *paymentDao) FindCreditSummariesWithPage(db *gorm.DB, ownerUser string, validStatuses []int, offset, limit int) ([]CreditSummary, int64, error) {
	var summaries []CreditSummary
	var total int64

	query := db.Model(&model.BatchOrder{}).
		Select("user_uuid, SUM(credit_amount) as total_credit").
		Where("owner_user = ? AND status IN (?) AND credit_amount != 0", ownerUser, validStatuses).
		Group("user_uuid")

	// 获取总数
	query.Count(&total)

	// 分页
	if limit > 0 {
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Scan(&summaries).Error
	return summaries, total, err
}

// FindOrdersWithCredit 查询客户所有未还清的订单，按时间排序
func (dao *paymentDao) FindOrdersWithCredit(db *gorm.DB, ownerUser, customerUUID string, validStatuses []int) ([]model.BatchOrder, error) {
	var orders []model.BatchOrder
	err := db.
		Where("owner_user = ? AND user_uuid = ? AND status IN (?) AND credit_amount != 0",
						ownerUser, customerUUID, validStatuses).
		Order("created_at ASC"). // 最早的订单优先
		Find(&orders).Error
	return orders, err
}

// FindOrderByUID 查询订单
func (dao *paymentDao) FindOrderByUID(db *gorm.DB, ownerUser, orderUUID string) (*model.BatchOrder, error) {
	var order model.BatchOrder
	err := db.Where("uid = ? AND owner_user = ?", orderUUID, ownerUser).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// CreatePay 创建还款记录（分表）
func (dao *paymentDao) CreatePay(db *gorm.DB, pay *model.BatchOrderPay) error {
	tableName := dao.getShardPayTable(pay.OwnerUser)
	return db.Table(tableName).Create(pay).Error
}

// FindPayRecords 查询还款记录
func (dao *paymentDao) FindPayRecords(db *gorm.DB, ownerUser, customerUUID string, offset, limit int) ([]model.BatchOrderPay, int64, error) {
	var pays []model.BatchOrderPay
	var total int64

	query := db.Model(&model.BatchOrderPay{}).
		Where("owner_user = ? AND customer_uuid = ?", ownerUser, customerUUID).
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
func (dao *paymentDao) FindPayByUID(db *gorm.DB, ownerUser, payUUID string) (*model.BatchOrderPay, error) {
	var pay model.BatchOrderPay
	err := db.Where("uid = ? AND owner_user = ?", payUUID, ownerUser).First(&pay).Error
	if err != nil {
		return nil, err
	}
	return &pay, nil
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

// getShardPayTable 获取分表名
func (dao *paymentDao) getShardPayTable(ownerUser string) string {
	// CRC32 hash 分表
	var hash uint32
	for _, c := range ownerUser {
		hash = hash*31 + uint32(c)
	}
	shardIndex := hash % 50
	return fmt.Sprintf("batch_order_pay_%d", shardIndex)
}
