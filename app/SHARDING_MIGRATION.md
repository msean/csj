# 分表DAO迁移指南

## 📋 DAO方法对照表

### OrderDao → ShardOrderDao

| 原方法 | 分表方法 | 变化 |
|--------|---------|------|
| `OrderDao.Create()` | `shardOrderDao.Create()` | ✅ 已实现 |
| `OrderDao.UpdateStatus()` | `shardOrderDao.UpdateStatus()` | ✅ 需要ownerUser参数 |
| `OrderDao.Shared()` | `shardOrderDao.Shared()` | ✅ 需要ownerUser参数 |
| `OrderDao.FindByUID()` | `shardOrderDao.FindByUID()` | ✅ 需要ownerUser参数 |
| `OrderDao.CreateGoods()` | `shardOrderDao.CreateGoods()` | ✅ 已实现 |
| `OrderDao.CreateGoodsBatch()` | `shardOrderDao.CreateGoodsBatch()` | ✅ 已实现 |
| `OrderDao.DeleteGoodsByOrderUUID()` | `shardOrderDao.DeleteGoodsByOrderUUID()` | ✅ 需要ownerUser参数 |
| `OrderDao.CustomerCreditAmount()` | `shardOrderDao.CustomerCreditAmount()` | ✅ 需要ownerUser参数 |
| `OrderDao.CreditAmountTotal()` | `shardOrderDao.CreditAmountTotal()` | ✅ 需要ownerUser参数 |
| `OrderDao.MonthSales()` | `shardOrderDao.MonthSales()` | ✅ 需要ownerUser参数 |
| `OrderDao.LatestOrderByCustomers()` | `shardOrderDao.LatestOrderByCustomers()` | ✅ 新增，需要ownerUser参数 |
| `OrderDao.ListByBatchUUIDIn()` | `shardOrderDao.ListByBatchUUIDIn()` | ✅ 新增，需要ownerUser参数 |
| `OrderDao.UpdateOrderPay()` | `shardOrderDao.UpdateOrderPay()` | ✅ 新增，需要ownerUser参数 |
| `OrderDao.UpdateGoods()` | `shardOrderDao.UpdateGoods()` | ✅ 新增，需要ownerUser参数 |
| `OrderDao.GetTodayOrdersWithGoods()` | `shardOrderDao.GetTodayOrdersWithGoods()` | ✅ 新增，需要ownerUser参数 |
| `OrderDao.GetCreditList()` | `shardOrderDao.GetCreditList()` | ✅ 新增，需要ownerUser参数 |

### BatchGoodsDao → ShardBatchGoodsDao

| 原方法 | 分表方法 | 变化 |
|--------|---------|------|
| `batchGoodsDao.FromUUID()` | `shardBatchGoodsDao.FromUUID()` | ✅ 需要ownerUser参数 |
| `batchGoodsDao.ListByBatchUUID()` | `shardBatchGoodsDao.ListByBatchUUID()` | ✅ 需要ownerUser参数 |
| - | `shardBatchGoodsDao.Create()` | ✅ 新增 |
| - | `shardBatchGoodsDao.CreateBatch()` | ✅ 新增 |
| - | `shardBatchGoodsDao.Update()` | ✅ 新增 |
| - | `shardBatchGoodsDao.DeleteByBatchUUID()` | ✅ 新增 |
| - | `shardBatchGoodsDao.DeleteByGoodsUUID()` | ✅ 新增 |

## 📊 分表结构

### 分表总数：150张

| 表类型 | 表名前缀 | 数量 | 说明 |
|--------|---------|------|------|
| 订单表 | `batch_orders_*` | 50 | 订单主表 |
| 订单货品表 | `batch_order_goods_*` | 50 | 订单货品明细 |
| 订单支付表 | `batch_order_pay_*` | 50 | 订单支付记录 |

### BatchOrderPay 分表结构

```sql
CREATE TABLE batch_order_pay_XXX (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    uid VARCHAR(64) NOT NULL UNIQUE,
    owner_user VARCHAR(64) NOT NULL,
    batch_order_uuid VARCHAR(64) NOT NULL,
    pay_type INT DEFAULT 0,
    amount DECIMAL(10,2) DEFAULT 0.00,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    INDEX idx_owner_user (owner_user),
    INDEX idx_batch_order_uuid (batch_order_uuid),
    INDEX idx_uid (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 🔧 代码迁移示例

### 1. customer.go - LatestOrderByCustomers

**迁移前**：
```go
bill, _ := dao.OrderDao.LatestOrderByCustomers(logic.runtime.DB, logic.OwnerUser, _customerModel)
```

**迁移后**：
```go
shardOrderDao := dao.newShardOrderDao()
bill, _ := shardOrderDao.LatestOrderByCustomers(logic.runtime.DB, logic.OwnerUser, _customerModel)
```

### 2. batch.go - ListByBatchUUIDIn

**迁移前**：
```go
orderGoodsList, err = dao.OrderDao.ListByBatchUUIDIn(logic.runtime.DB, logic.UID, toDeleteGoodsUUIDs)
```

**迁移后**：
```go
shardOrderDao := dao.newShardOrderDao()
orderGoodsList, err = shardOrderDao.ListByBatchUUIDIn(logic.runtime.DB, logic.OwnerUser, logic.UID, toDeleteGoodsUUIDs)
```

### 3. order_pay.go - UpdateOrderPay

**迁移前**：
```go
if err = dao.OrderDao.UpdateOrderPay(tx, logic.BatchOrderPay); err != nil {
    return
}
```

**迁移后**：
```go
shardOrderDao := dao.newShardOrderDao()
if err = shardOrderDao.UpdateOrderPay(tx, logic.OwnerUser, logic.BatchOrderPay); err != nil {
    return
}
```

### 3.1 order_pay.go - CreatePay

**迁移前**：
```go
if err = utils.CreateObj(tx, &logic.BatchOrderPay); err != nil {
    return
}
```

**迁移后**：
```go
shardOrderDao := dao.newShardOrderDao()
if err = shardOrderDao.CreatePay(tx, &logic.BatchOrderPay); err != nil {
    return
}
```

### 4. order.go - GetTodayOrdersWithGoods

**迁移前**：
```go
orders, previousCredit, err := dao.OrderDao.GetTodayOrdersWithGoods(logic.runtime.DB, logic.OwnerUser, req.CustomerUUID)
```

**迁移后**：
```go
shardOrderDao := dao.newShardOrderDao()
orders, previousCredit, err := shardOrderDao.GetTodayOrdersWithGoods(logic.runtime.DB, logic.OwnerUser, req.CustomerUUID)
```

## ⚠️ 重要规则

### 1. 所有分表DAO方法都需要 ownerUser 参数

```go
// ❌ 错误：缺少ownerUser
order, err := shardOrderDao.FindByUID(db, uid)

// ✅ 正确：携带ownerUser
order, err := shardOrderDao.FindByUID(db, ownerUser, uid)
```

### 2. 初始化分表DAO实例

建议在Logic层缓存DAO实例：

```go
type OrderLogic struct {
    context       *gin.Context
    runtime       *global.RunTime
    model.BatchOrder
    shardOrderDao *dao.shardOrderDao  // 缓存实例
}

func NewOrderLogic(context *gin.Context) *OrderLogic {
    logic := &OrderLogic{
        context:       context,
        runtime:       global.Global,
        shardOrderDao: dao.newShardOrderDao(),  // 初始化
    }
    logic.OwnerUser = common.GetUserUUID(context)
    return logic
}
```

### 3. 事务中使用分表DAO

```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// 所有操作使用同一个tx
shardOrderDao.Create(tx, &order)
shardOrderDao.CreateGoodsBatch(tx, order.OwnerUser, goodsList)

tx.Commit()
```

## 📝 待迁移的文件清单

根据您提供的代码，以下文件需要迁移：

- [ ] `app/service/logic/customer.go` - Line 70
- [ ] `app/service/logic/batch.go` - Line 174
- [ ] `app/service/logic/order_pay.go` - Line 56
- [ ] `app/service/logic/order.go` - Line 249
- [ ] `app/service/logic/order.go` - 其他OrderDao调用
- [ ] `app/service/logic/batch.go` - 其他BatchGoodsDao调用

## 🚀 快速迁移脚本

如果您需要批量替换，可以使用以下模式：

```bash
# 查找所有需要替换的地方
grep -r "dao\.OrderDao\." app/service/logic/
grep -r "dao\.BatchGoodsDao\." app/service/logic/
```

## ✅ 迁移检查清单

- [ ] 所有DAO调用都传入了ownerUser
- [ ] 初始化了shardOrderDao实例
- [ ] 初始化了shardBatchGoodsDao实例  
- [ ] 事务操作正确使用分表DAO
- [ ] 预加载(GoodsListRelated)使用分表版本
- [ ] 统计查询使用分表版本
- [ ] 测试通过
