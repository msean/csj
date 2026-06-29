package model

import (
	"app/global"
	"app/pkg/utils"
	"fmt"
	"log"

	"go.uber.org/zap"
)

// MigrateShardingTables 创建50张分表
func MigrateShardingTables() error {
	global.Global.Logger.Info("开始创建分表...")

	// 0. 删除外键约束（分表架构不需要）
	if err := dropForeignKeyConstraints(); err != nil {
		global.Global.Logger.Warn("删除外键约束失败（可能已删除）", zap.Error(err))
	}

	// 1. 创建 batch_orders 分表
	orderTables := utils.GetAllShardTableNames("batch_orders")
	for i, tableName := range orderTables {
		if err := createOrderTable(tableName); err != nil {
			return fmt.Errorf("创建 batch_orders_%d 失败: %w", i, err)
		}
		global.Global.Logger.Info(fmt.Sprintf("创建分表成功: %s", tableName))
	}

	// 创建 batch_order_goods 分表
	goodsTables := utils.GetAllShardTableNames("batch_order_goods")
	for i, tableName := range goodsTables {
		if err := createGoodsTable(tableName); err != nil {
			return fmt.Errorf("创建 batch_order_goods_%d 失败: %w", i, err)
		}
		global.Global.Logger.Info(fmt.Sprintf("创建分表成功: %s", tableName))
	}

	// 创建 batch_order_pay 分表
	payTables := utils.GetAllShardTableNames("batch_order_pay")
	for i, tableName := range payTables {
		if err := createPayTable(tableName); err != nil {
			return fmt.Errorf("创建 batch_order_pay_%d 失败: %w", i, err)
		}
		global.Global.Logger.Info(fmt.Sprintf("创建分表成功: %s", tableName))
	}

	// 创建 batch_goods 分表
	return nil
}

// MigrateShardingData 将现有数据迁移到分表中
func MigrateShardingData() error {
	global.Global.Logger.Info("开始迁移数据到分表...")

	// 0. 修复已创建的batch_order_pay分表（添加customer_uuid字段）
	if err := fixPayTableStructure(); err != nil {
		return fmt.Errorf("修复支付表结构失败: %w", err)
	}

	// 1. 迁移 batch_orders 数据
	if err := migrateOrdersData(); err != nil {
		return fmt.Errorf("迁移订单数据失败: %w", err)
	}

	// 2. 迁移 batch_order_goods 数据
	if err := migrateOrderGoodsData(); err != nil {
		return fmt.Errorf("迁移订单货品数据失败: %w", err)
	}

	// 3. 迁移 batch_order_pay 数据
	if err := migrateOrderPayData(); err != nil {
		return fmt.Errorf("迁移订单支付数据失败: %w", err)
	}

	global.Global.Logger.Info("数据迁移完成！")
	return nil
}

// fixPayTableStructure 修复已创建的batch_order_pay分表结构
func fixPayTableStructure() error {
	log.Println("[FixPayTable] 开始修复batch_order_pay分表结构...")

	payTables := utils.GetAllShardTableNames("batch_order_pay")
	fixedCount := 0

	for _, tableName := range payTables {
		// 检查是否已有customer_uuid字段
		var colCount int64
		err := global.Global.DB.Raw(
			"SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? AND COLUMN_NAME = 'customer_uuid'",
			tableName).Scan(&colCount).Error

		if err != nil {
			log.Printf("[FixPayTable] 检查表 %s 结构失败: %v", tableName, err)
			continue
		}

		if colCount > 0 {
			continue // 字段已存在，跳过
		}

		// 检查表是否存在
		var tableCount int64
		global.Global.DB.Raw(
			"SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?",
			tableName).Scan(&tableCount)

		if tableCount == 0 {
			continue // 表还未创建，跳过
		}

		// 添加customer_uuid字段
		sql := fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN customer_uuid VARCHAR(64) NOT NULL DEFAULT '' AFTER uid",
			tableName)

		if err := global.Global.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("修复表 %s 失败: %w", tableName, err)
		}

		// 添加索引
		indexSQL := fmt.Sprintf(
			"ALTER TABLE %s ADD INDEX idx_customer_uuid (customer_uuid)",
			tableName)
		global.Global.DB.Exec(indexSQL) // 忽略错误，可能索引已存在

		fixedCount++
		log.Printf("[FixPayTable] 修复表 %s 成功", tableName)
	}

	if fixedCount > 0 {
		log.Printf("[FixPayTable] 共修复了 %d 张batch_order_pay分表", fixedCount)
	} else {
		log.Println("[FixPayTable] 所有batch_order_pay分表结构正确")
	}

	return nil
}

// migrateOrdersData 迁移订单数据
func migrateOrdersData() error {
	log.Println("[MigrateOrders] 开始迁移订单数据...")

	// 查询所有订单数据
	var orders []BatchOrder
	if err := global.Global.DB.Find(&orders).Error; err != nil {
		return fmt.Errorf("查询订单数据失败: %w", err)
	}

	if len(orders) == 0 {
		log.Println("[MigrateOrders] 没有订单数据需要迁移")
		return nil
	}

	log.Printf("[MigrateOrders] 找到 %d 条订单数据", len(orders))

	// 按 owner_user 分组
	orderMap := make(map[string][]BatchOrder)
	for _, order := range orders {
		orderMap[order.OwnerUser] = append(orderMap[order.OwnerUser], order)
	}

	// 迁移到分表
	totalMigrated := 0
	for ownerUser, userOrders := range orderMap {
		shardTable := utils.GetShardTableName("batch_orders", ownerUser)

		// 检查是否已迁移（避免重复）
		var count int64
		global.Global.DB.Table(shardTable).Count(&count)
		if count > 0 {
			log.Printf("[MigrateOrders] 分表 %s 已有 %d 条数据，跳过", shardTable, count)
			continue
		}

		// 批量插入
		if err := global.Global.DB.Table(shardTable).CreateInBatches(userOrders, 100).Error; err != nil {
			return fmt.Errorf("插入分表 %s 失败: %w", shardTable, err)
		}

		totalMigrated += len(userOrders)
		log.Printf("[MigrateOrders] 迁移 %d 条订单到 %s", len(userOrders), shardTable)
	}

	log.Printf("[MigrateOrders] 订单数据迁移完成，共迁移 %d 条", totalMigrated)
	return nil
}

// migrateOrderGoodsData 迁移订单货品数据
func migrateOrderGoodsData() error {
	log.Println("[MigrateOrderGoods] 开始迁移订单货品数据...")

	// 查询所有订单货品数据
	var goods []BatchOrderGoods
	if err := global.Global.DB.Find(&goods).Error; err != nil {
		return fmt.Errorf("查询订单货品数据失败: %w", err)
	}

	if len(goods) == 0 {
		log.Println("[MigrateOrderGoods] 没有订单货品数据需要迁移")
		return nil
	}

	log.Printf("[MigrateOrderGoods] 找到 %d 条订单货品数据", len(goods))

	// 按 owner_user 分组
	goodsMap := make(map[string][]BatchOrderGoods)
	for _, good := range goods {
		goodsMap[good.OwnerUser] = append(goodsMap[good.OwnerUser], good)
	}

	// 迁移到分表
	totalMigrated := 0
	for ownerUser, userGoods := range goodsMap {
		shardTable := utils.GetShardTableName("batch_order_goods", ownerUser)

		// 检查是否已迁移
		var count int64
		global.Global.DB.Table(shardTable).Count(&count)
		if count > 0 {
			log.Printf("[MigrateOrderGoods] 分表 %s 已有 %d 条数据，跳过", shardTable, count)
			continue
		}

		// 批量插入
		if err := global.Global.DB.Table(shardTable).CreateInBatches(userGoods, 100).Error; err != nil {
			return fmt.Errorf("插入分表 %s 失败: %w", shardTable, err)
		}

		totalMigrated += len(userGoods)
		log.Printf("[MigrateOrderGoods] 迁移 %d 条订单货品到 %s", len(userGoods), shardTable)
	}

	log.Printf("[MigrateOrderGoods] 订单货品数据迁移完成，共迁移 %d 条", totalMigrated)
	return nil
}

// migrateOrderPayData 迁移订单支付数据
func migrateOrderPayData() error {
	log.Println("[MigrateOrderPay] 开始迁移订单支付数据...")

	// 查询所有订单支付数据
	var pays []BatchOrderPay
	if err := global.Global.DB.Find(&pays).Error; err != nil {
		return fmt.Errorf("查询订单支付数据失败: %w", err)
	}

	if len(pays) == 0 {
		log.Println("[MigrateOrderPay] 没有订单支付数据需要迁移")
		return nil
	}

	log.Printf("[MigrateOrderPay] 找到 %d 条订单支付数据", len(pays))

	// 按 owner_user 分组
	payMap := make(map[string][]BatchOrderPay)
	for _, pay := range pays {
		payMap[pay.OwnerUser] = append(payMap[pay.OwnerUser], pay)
	}

	// 迁移到分表
	totalMigrated := 0
	for ownerUser, userPays := range payMap {
		shardTable := utils.GetShardTableName("batch_order_pay", ownerUser)

		// 检查是否已迁移
		var count int64
		global.Global.DB.Table(shardTable).Count(&count)
		if count > 0 {
			log.Printf("[MigrateOrderPay] 分表 %s 已有 %d 条数据，跳过", shardTable, count)
			continue
		}

		// 批量插入
		if err := global.Global.DB.Table(shardTable).CreateInBatches(userPays, 100).Error; err != nil {
			return fmt.Errorf("插入分表 %s 失败: %w", shardTable, err)
		}

		totalMigrated += len(userPays)
		log.Printf("[MigrateOrderPay] 迁移 %d 条订单支付到 %s", len(userPays), shardTable)
	}

	log.Printf("[MigrateOrderPay] 订单支付数据迁移完成，共迁移 %d 条", totalMigrated)
	return nil
}

// createOrderTable 创建订单分表
func createOrderTable(tableName string) error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			uid VARCHAR(64) NOT NULL PRIMARY KEY,
			batch_uuid VARCHAR(64) NOT NULL,
			owner_user VARCHAR(64) NOT NULL,
			user_uuid VARCHAR(64) NOT NULL,
			shared INT DEFAULT 0,
			share_time INT DEFAULT 0,
			status INT DEFAULT 0,
			amount DECIMAL(10,2) DEFAULT 0.00,
			credit_amount DECIMAL(10,2) DEFAULT 0.00,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			INDEX idx_owner_user (owner_user),
			INDEX idx_user_uuid (user_uuid),
			INDEX idx_batch_uuid (batch_uuid),
			INDEX idx_created_at (created_at),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`, tableName)
	return global.Global.DB.Exec(sql).Error
}

// createGoodsTable 创建订单货品分表
func createGoodsTable(tableName string) error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			uid VARCHAR(64) NOT NULL PRIMARY KEY,
			batch_order_uuid VARCHAR(64) NOT NULL,
			batch_uuid VARCHAR(64) NOT NULL,
			goods_uuid VARCHAR(64) NOT NULL,
			owner_user VARCHAR(64) NOT NULL,
			user_uuid VARCHAR(64) NOT NULL,
			good_type INT DEFAULT 0,
			price DECIMAL(10,2) DEFAULT 0.00,
			weight DECIMAL(10,2) DEFAULT 0.00,
			mount INT DEFAULT 0,
			total DECIMAL(10,2) DEFAULT 0.00,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			INDEX idx_owner_user (owner_user),
			INDEX idx_batch_order_uuid (batch_order_uuid),
			INDEX idx_batch_uuid (batch_uuid),
			INDEX idx_goods_uuid (goods_uuid),
			INDEX idx_user_uuid (user_uuid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`, tableName)
	return global.Global.DB.Exec(sql).Error
}

// createPayTable 创建订单支付分表
func createPayTable(tableName string) error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			uid VARCHAR(64) NOT NULL PRIMARY KEY,
			customer_uuid VARCHAR(64) NOT NULL,
			batch_order_uuid VARCHAR(64) NOT NULL,
			owner_user VARCHAR(64) NOT NULL,
			pay_type INT DEFAULT 0,
			amount DECIMAL(10,2) DEFAULT 0.00,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			deleted_at DATETIME,
			INDEX idx_owner_user (owner_user),
			INDEX idx_batch_order_uuid (batch_order_uuid),
			INDEX idx_customer_uuid (customer_uuid),
			INDEX idx_uid (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`, tableName)
	return global.Global.DB.Exec(sql).Error
}

// dropForeignKeyConstraints 删除外键约束（分表架构不需要）
func dropForeignKeyConstraints() error {
	log.Println("[DropForeignKey] 开始删除外键约束...")

	// 删除 batch_order_goods 表的外键约束
	fkName := "fk_batch_orders_goods_list_related"
	tableName := "batch_order_goods"

	// 检查外键是否存在
	var count int64
	global.Global.DB.Raw(
		"SELECT COUNT(*) FROM information_schema.TABLE_CONSTRAINTS WHERE CONSTRAINT_SCHEMA = DATABASE() AND CONSTRAINT_NAME = ? AND TABLE_NAME = ?",
		fkName, tableName).Scan(&count)

	if count > 0 {
		sql := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", tableName, fkName)
		if err := global.Global.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("删除外键 %s 失败: %w", fkName, err)
		}
		log.Printf("[DropForeignKey] 删除外键约束 %s 成功", fkName)
	} else {
		log.Printf("[DropForeignKey] 外键约束 %s 不存在，跳过", fkName)
	}

	return nil
}
