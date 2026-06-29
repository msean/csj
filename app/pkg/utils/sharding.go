package utils

import (
	"fmt"
	"hash/crc32"
)

const (
	// ShardingTableCount 分表数量
	ShardingTableCount = 50
)

// GetShardIndex 根据OwnerUser计算分表索引 (0-49)
func GetShardIndex(ownerUser string) int {
	if ownerUser == "" {
		return 0
	}
	// 使用CRC32哈希算法
	hash := crc32.ChecksumIEEE([]byte(ownerUser))
	return int(hash) % ShardingTableCount
}

// GetShardTableName 获取分表名称
// tableName: 基础表名 (如 "batch_orders")
// ownerUser: 用户标识
// 返回: "batch_orders_0", "batch_orders_1", ... "batch_orders_49"
func GetShardTableName(tableName string, ownerUser string) string {
	shardIndex := GetShardIndex(ownerUser)
	return fmt.Sprintf("%s_%d", tableName, shardIndex)
}

// GetAllShardTableNames 获取所有分表名称
func GetAllShardTableNames(tableName string) []string {
	tables := make([]string, ShardingTableCount)
	for i := 0; i < ShardingTableCount; i++ {
		tables[i] = fmt.Sprintf("%s_%d", tableName, i)
	}
	return tables
}
