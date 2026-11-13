package dao

import (
	"app/service/common"
	"fmt"
	"gorm.io/gorm"
)

type BatchInventoryUsage struct {
	GoodsUUID string
	Mount     int
	Weight    float64
}

// BatchInventoryUsageQuery aggregates only sold quantities, with the same valid
// order statuses as full batch statistics. No prices, revenue or profit are read.
func BatchInventoryUsageQuery(db *gorm.DB, owner, batch string) *gorm.DB {
	return db.Table(fmt.Sprintf("%s AS bog", ShardOrderGoodsTable(owner))).
		Select("bog.goods_uuid, COALESCE(SUM(bog.mount), 0) AS mount, COALESCE(SUM(bog.weight), 0) AS weight").
		Joins(fmt.Sprintf("JOIN %s AS bo ON bo.uid = bog.batch_order_uuid", ShardOrderTable(owner))).
		Where("bog.owner_user = ? AND bo.owner_user = ? AND bog.batch_uuid = ?", owner, owner, batch).
		Where("bo.status IN (?)", common.ValidOrder).
		Where("bog.deleted_at IS NULL AND bo.deleted_at IS NULL").
		Group("bog.goods_uuid")
}
