package dao

import (
	"app/service/model"

	"gorm.io/gorm"
)

func CustomerPurchaseQuery(db *gorm.DB, owner, customer string, start, end int64, creditOnly bool) *gorm.DB {
	q := db.Table(ShardOrderTable(owner)).Where("owner_user = ? AND user_uuid = ? AND deleted_at IS NULL AND status IN (1, 2, 3)", owner, customer)
	if start > 0 {
		q = q.Where("order_time >= ?", start)
	}
	if end > 0 {
		q = q.Where("order_time < ?", end)
	}
	if creditOnly {
		q = q.Where("credit_amount > 0")
	}
	return q
}

func CustomerPurchaseSummary(db *gorm.DB, owner, customer string, start, end int64, creditOnly bool) (total int64, amount float64, err error) {
	var result struct {
		Total  int64
		Amount float64
	}
	err = CustomerPurchaseQuery(db, owner, customer, start, end, creditOnly).
		Select("COUNT(*) AS total, COALESCE(SUM(amount + COALESCE(tricycle_fee, 0)), 0) AS amount").Scan(&result).Error
	return result.Total, result.Amount, err
}

func CustomerPurchaseList(db *gorm.DB, owner, customer string, start, end int64, creditOnly bool, offset, limit int) ([]*model.BatchOrder, error) {
	rows := make([]*model.BatchOrder, 0)
	err := CustomerPurchaseQuery(db, owner, customer, start, end, creditOnly).
		Preload("GoodsListRelated", func(q *gorm.DB) *gorm.DB {
			return OrderGoodsInSequence(q.Table(ShardOrderGoodsTable(owner)).Where("owner_user = ?", owner))
		}).
		Order("order_time DESC, uid DESC").Offset(offset).Limit(limit).Find(&rows).Error
	return rows, err
}
