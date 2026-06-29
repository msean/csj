# BatchOrder 分表方案使用指南

## 📊 分表策略

- **分表数量**: 50张表
- **分表算法**: `CRC32(OwnerUser) % 50`
- **表名规则**: 
  - 订单表: `batch_orders_0` ~ `batch_orders_49`
  - 货品表: `batch_order_goods_0` ~ `batch_order_goods_49`

## 🚀 快速开始

### 1. 初始化分表

在应用启动时调用（仅执行一次）：

```go
import "app/service/model"

// 创建50张分表
if err := model.MigrateShardingTables(); err != nil {
    log.Fatal("分表创建失败:", err)
}
```

### 2. 使用分表DAO

```go
import "app/service/dao"

// 获取分表DAO实例
shardOrderDao := dao.newShardOrderDao()

// 创建订单（自动路由到对应分表）
err := shardOrderDao.Create(db, &order)

// 查询订单（需要传入ownerUser）
order, err := shardOrderDao.FindByUID(db, ownerUser, uid, "GoodsListRelated")

// 更新状态
err := shardOrderDao.UpdateStatus(db, ownerUser, uid, status)
```

### 3. 手动计算分表索引

```go
import "app/pkg/utils"

// 计算分表索引 (0-49)
index := utils.GetShardIndex(ownerUser)

// 获取分表名
tableName := utils.GetShardTableName("batch_orders", ownerUser)
// 返回: "batch_orders_23" (示例)

// 获取所有分表名
allTables := utils.GetAllShardTableNames("batch_orders")
// 返回: ["batch_orders_0", "batch_orders_1", ..., "batch_orders_49"]
```

## ⚠️ 重要注意事项

### 1. 查询必须携带 OwnerUser

分表查询**必须**知道 `ownerUser`，否则无法路由到正确的分表：

```go
// ✅ 正确：携带 ownerUser
order, err := shardOrderDao.FindByUID(db, ownerUser, uid)

// ❌ 错误：没有 ownerUser，无法确定分表
order, err := shardOrderDao.FindByUID(db, "", uid) // 会路由到 batch_orders_0
```

### 2. 跨分表查询

如果需要查询所有用户的数据（如管理员统计），需要遍历所有50张表：

```go
// 示例：统计所有用户的订单总数
var totalCount int64
allTables := utils.GetAllShardTableNames("batch_orders")

for _, tableName := range allTables {
    var count int64
    db.Table(tableName).Count(&count)
    totalCount += count
}
```

### 3. 事务支持

同一个订单和它的货品会在同一组分表中（因为使用相同的OwnerUser）：

```go
tx := db.Begin()

// 创建订单
err := shardOrderDao.Create(tx, &order)

// 创建货品（自动路由到对应的 goods 分表）
err = shardOrderDao.CreateGoodsBatch(tx, order.OwnerUser, goodsList)

tx.Commit()
```

### 4. 数据迁移

如果已有数据在旧表中，需要迁移到分表：

```sql
-- 示例：将 batch_orders 数据迁移到分表
-- 需要在应用层编写迁移脚本，根据 owner_user 计算分表索引
```

## 📈 性能优势

- **查询性能**: 单表数据量减少到 1/50
- **索引效率**: 每个分表独立索引，查询更快
- **并发写入**: 不同用户的写入操作分散到不同表
- **维护便捷**: 可以单独备份/优化某个分表

## 🔧 分表工具函数

位于 `app/pkg/utils/sharding.go`:

| 函数 | 说明 | 示例 |
|------|------|------|
| `GetShardIndex(ownerUser)` | 获取分表索引(0-49) | `GetShardIndex("user123") → 23` |
| `GetShardTableName(base, ownerUser)` | 获取分表名 | `GetShardTableName("batch_orders", "user123") → "batch_orders_23"` |
| `GetAllShardTableNames(base)` | 获取所有分表名 | `GetAllShardTableNames("batch_orders") → ["batch_orders_0", ...]` |

## 🎯 最佳实践

1. **始终携带 OwnerUser**: 所有查询/更新/删除操作都要传入 ownerUser
2. **批量操作**: 使用 `CreateGoodsBatch` 而不是循环插入
3. **避免跨表JOIN**: 订单和货品天然在同一分片，可以安全JOIN
4. **监控分表均衡**: 定期检查各分表数据量是否均匀分布
5. **备份策略**: 可以按分表备份，提高恢复效率

## ❓ FAQ

**Q: 如果 OwnerUser 为空会怎样？**  
A: 会路由到 `batch_orders_0`，建议始终确保 ownerUser 不为空。

**Q: 分表后可以改成100张吗？**  
A: 需要数据迁移。建议初期就规划好分表数量。

**Q: 如何统计全局数据？**  
A: 需要遍历50张表聚合，或使用独立的统计汇总表。

**Q: 分表后外键怎么办？**  
A: 分表不支持外键约束，需要在应用层保证数据一致性。
