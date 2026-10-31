package utils

import "github.com/msean/csj/shared/sharding"

const ShardingTableCount = sharding.ShardingTableCount

func GetShardIndex(owner string) int               { return sharding.GetShardIndex(owner) }
func GetShardTableName(table, owner string) string { return sharding.GetShardTableName(table, owner) }
func GetAllShardTableNames(table string) []string  { return sharding.GetAllShardTableNames(table) }
