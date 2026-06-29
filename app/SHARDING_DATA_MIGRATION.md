# 分表数据迁移指南

## 📋 概述

本次迁移会将现有数据从单表迁移到50个分表中，支持以下表：
- `batch_orders` → `batch_orders_0` ~ `batch_orders_49`
- `batch_order_goods` → `batch_order_goods_0` ~ `batch_order_goods_49`
- `batch_order_pay` → `batch_order_pay_0` ~ `batch_order_pay_49`

## 🚀 自动迁移（推荐）

系统启动时会自动执行迁移，无需手动操作：

```bash
# 启动服务
./app -c config.yaml
```

启动日志会显示：
```
[INFO] 开始创建分表...
[INFO] 创建分表成功: batch_orders_0
[INFO] 创建分表成功: batch_order_goods_0
[INFO] 创建分表成功: batch_order_pay_0
...
[INFO] 开始迁移数据到分表...
[MigrateOrders] 开始迁移订单数据...
[MigrateOrders] 找到 1234 条订单数据
[MigrateOrders] 迁移 1234 条订单到 batch_orders_12
[MigrateOrders] 订单数据迁移完成，共迁移 1234 条
...
[INFO] 数据迁移完成！
```

## ✅ 迁移特性

### 1. **幂等性**
迁移脚本是幂等的，可以安全地多次执行：
- ✅ 如果分表已存在数据，会跳过该分表
- ✅ 不会重复迁移数据
- ✅ 不会丢失数据

### 2. **按用户分组**
数据按 `owner_user` 分组迁移到对应分表：
```
owner_user = "abc123" → CRC32("abc123") % 50 = 12
所有数据迁移到 *_12 分表
```

### 3. **批量插入**
使用批量插入提高性能（每批100条）：
```go
db.Table(shardTable).CreateInBatches(data, 100)
```

## 🔍 手动验证迁移

### 1. 检查分表是否创建

```sql
-- 查看所有分表
SHOW TABLES LIKE 'batch_orders_%';
SHOW TABLES LIKE 'batch_order_goods_%';
SHOW TABLES LIKE 'batch_order_pay_%';

-- 应该各有50张表
```

### 2. 检查数据是否正确迁移

```sql
-- 对比原表和分表的数据量
SELECT COUNT(*) FROM batch_orders;
SELECT SUM(TABLE_ROWS) FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'your_database' 
AND TABLE_NAME LIKE 'batch_orders_%';

-- 两个数字应该相同
```

### 3. 检查特定用户的数据

```sql
-- 查询某个用户的数据在哪个分表
SELECT owner_user, COUNT(*) 
FROM batch_orders 
GROUP BY owner_user;

-- 验证分表数据
SELECT * FROM batch_orders_12 LIMIT 10;
```

## ⚠️ 注意事项

### 1. **首次执行**
- 迁移仅在首次启动时执行
- 如果分表已有数据，会跳过
- 建议在低峰期执行首次迁移

### 2. **数据一致性**
- 迁移过程中，原表数据保持不变
- 迁移完成后，需要切换业务代码使用分表DAO
- 确认分表数据正确后，可以删除原表（谨慎！）

### 3. **性能影响**
- 数据量小时（<10万条），迁移几乎无影响
- 数据量大时，建议在维护窗口执行
- 迁移过程会占用数据库资源

## 🛠️ 故障排查

### 问题1：迁移失败

**症状**：启动时 panic，日志显示迁移错误

**解决**：
```bash
# 1. 检查数据库连接
mysql -u root -p -e "SHOW DATABASES;"

# 2. 检查表结构
SHOW CREATE TABLE batch_orders;

# 3. 查看详细日志
tail -f app.log | grep -i migrate
```

### 问题2：分表数据不完整

**症状**：分表数据量与原表不一致

**解决**：
```sql
-- 1. 检查哪些分表有数据
SELECT TABLE_NAME, TABLE_ROWS 
FROM information_schema.TABLES 
WHERE TABLE_NAME LIKE 'batch_orders_%';

-- 2. 重新执行迁移（安全，会跳过已有数据）
./app -c config.yaml
```

### 问题3：重复迁移

**症状**：日志显示"分表已有X条数据，跳过"

**说明**：这是正常行为，说明迁移已执行过，数据已存在。

## 📊 迁移后验证清单

- [ ] 所有150张分表已创建
- [ ] 原表数据量 = 分表数据量总和
- [ ] 随机抽查用户数据在正确分表
- [ ] 业务代码已切换使用分表DAO
- [ ] 新数据写入分表（而非原表）
- [ ] 查询性能正常

## 🎯 迁移后下一步

1. **切换业务代码**：将所有 DAO 调用改为分表版本
2. **测试验证**：确保所有功能正常
3. **监控性能**：观察分表后的性能表现
4. **清理原表**：确认无误后，备份并删除原表（可选）

## 📞 需要帮助？

查看迁移日志：
```bash
tail -100 app.log | grep -E "Migrate|分表"
```

查看数据库状态：
```sql
SELECT 
    TABLE_NAME,
    TABLE_ROWS,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'your_database'
AND TABLE_NAME LIKE 'batch_%'
ORDER BY TABLE_NAME;
```
