# 收款功能实现文档

## 📋 功能概述

完整的收款管理系统，包括快捷还款、订单还款、还款历史、撤销还款和消息中心。

---

## 🗄️ 数据模型

### 1. BatchOrderPay 还款记录表

```go
type BatchOrderPay struct {
    BaseModel
    CustomerUUID   string     // 客户UUID
    OwnerUser      string     // 所属用户
    BatchOrderUUID string     // 订单UUID（快捷还款为空）
    Amount         float64    // 还款金额（正数还款，负数撤销）
    PayType        int32      // 付款方式 1-现金 2-微信 3-支付宝
    Remark         string     // 备注
    IsRevoked      int        // 是否已撤销 0-否 1-是
    RevokedAt      *time.Time // 撤销时间
    RevokedReason  string     // 撤销原因
    PayDetails     string     // 还款详情JSON（快捷还款使用）
}
```

**PayDetails JSON 格式**（快捷还款时存储）：
```json
[
  {"orderUUID": "xxx", "amount": 100},
  {"orderUUID": "yyy", "amount": 200}
]
```

### 2. MessageCenter 消息中心表

```go
type MessageCenter struct {
    BaseModel
    OwnerUser    string // 所属用户
    CustomerUUID string // 客户UUID
    Type         int    // 消息类型 1-还款 2-撤销还款
    Event        string // 事件名称
    Content      string // 消息内容
    IsRead       int    // 是否已读 0-未读 1-已读
    RelatedUUID  string // 关联UUID（还款记录UUID）
}
```

---

## 🔌 API 接口

### 1. 快速还款列表

**POST** `/api/csj/payment/quick/list`

**请求参数**：
```json
{
  "page": 1,
  "pageCount": 10
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "customerUUID": "uuid1",
        "customerName": "张三",
        "totalCredit": 900.00
      }
    ],
    "total": 1
  }
}
```

**说明**：显示 OwnerUser 下所有有赊欠的客户及其总赊欠金额。

---

### 2. 快捷还款

**POST** `/api/csj/payment/quick/pay`

**请求参数**：
```json
{
  "customerUUID": "uuid1",
  "amount": 600.00,
  "payType": 1,
  "remark": "现金还款"
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "payUUID": "pay-uuid-xxx"
  }
}
```

**还款规则**：
- 按订单创建时间**从早到晚**顺序还款
- 优先还清最早的订单
- 例如：
  - 订单A：赊欠100（最早）
  - 订单B：赊欠100
  - 订单C：赊欠300
  - 订单D：赊欠400（最晚）
  - 还款600后：
    - 订单A：赊欠0 ✅ 已还清
    - 订单B：赊欠0 ✅ 已还清
    - 订单C：赊欠100
    - 订单D：赊欠400

**注意**：`CreditAmount` 可以为负数，表示多付款了。

---

### 3. 针对订单还款

**POST** `/api/csj/payment/order/pay`

**请求参数**：
```json
{
  "orderUUID": "order-uuid-xxx",
  "amount": 200.00,
  "payType": 2,
  "remark": "微信还款"
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "payUUID": "pay-uuid-xxx",
    "orderUUID": "order-uuid-xxx",
    "orderCredit": -50.00  // 还款后订单赊欠金额（可以为负）
  }
}
```

**说明**：直接修改指定订单的 `CreditAmount` 字段。

---

### 4. 还款历史

**POST** `/api/csj/payment/history`

**请求参数**：
```json
{
  "customerUUID": "uuid1",
  "page": 1,
  "pageCount": 10
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "payUUID": "pay-uuid-xxx",
        "customerUUID": "uuid1",
        "customerName": "张三",
        "amount": 600.00,
        "payType": 1,
        "remark": "现金还款",
        "isRevoked": 0,
        "revokedAt": null,
        "revokedReason": "",
        "createdAt": "2026-06-17T10:30:00Z"
      }
    ],
    "total": 1
  }
}
```

**说明**：在客户详情页查看该客户的还款历史。

---

### 5. 还款详情

**POST** `/api/csj/payment/detail`

**请求参数**：
```json
{
  "payUUID": "pay-uuid-xxx"
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "payUUID": "pay-uuid-xxx",
    "customerUUID": "uuid1",
    "customerName": "张三",
    "orderUUID": "",  // 快捷还款为空
    "amount": 600.00,
    "payType": 1,
    "remark": "现金还款",
    "payDetails": "[{\"orderUUID\":\"order-a\",\"amount\":100},{\"orderUUID\":\"order-b\",\"amount\":100},{\"orderUUID\":\"order-c\",\"amount\":300},{\"orderUUID\":\"order-d\",\"amount\":100}]",
    "isRevoked": 0,
    "revokedAt": null,
    "revokedReason": "",
    "createdAt": "2026-06-17T10:30:00Z"
  }
}
```

---

### 6. 撤销还款

**POST** `/api/csj/payment/revoke`

**请求参数**：
```json
{
  "payUUID": "pay-uuid-xxx",
  "reason": "还款金额错误"
}
```

**响应**：
```json
{
  "code": 0,
  "data": null
}
```

**撤销规则**：
- **快捷还款**：按 `PayDetails` 中记录的每个订单和金额原路返回
- **订单还款**：直接恢复该订单的赊欠金额
- 已撤销的记录不能再次撤销
- 撤销后会在消息中心记录

---

### 7. 消息列表

**POST** `/api/csj/payment/message/list`

**请求参数**：
```json
{
  "page": 1,
  "pageCount": 20
}
```

**响应**：
```json
{
  "code": 0,
  "data": {
    "items": [
      {
        "messageUUID": "msg-uuid-xxx",
        "type": 1,  // 1-还款 2-撤销还款
        "event": "快捷还款",
        "content": "还款 600.00 元",
        "customerUUID": "uuid1",
        "customerName": "张三",
        "isRead": 0,
        "relatedUUID": "pay-uuid-xxx",
        "createdAt": "2026-06-17T10:30:00Z"
      }
    ],
    "summary": {
      "totalCount": 10,
      "unreadCount": 5
    }
  }
}
```

---

## 💡 核心逻辑

### 快捷还款算法

```go
// 1. 查询该客户所有未还清的订单，按时间排序
orders = SELECT * FROM batch_orders 
         WHERE user_uuid = ? AND credit_amount != 0
         ORDER BY created_at ASC

// 2. 按时间顺序还款
remainingAmount = req.Amount
for order in orders:
    if remainingAmount <= 0:
        break
    
    payAmount = min(remainingAmount, order.CreditAmount)
    order.CreditAmount -= payAmount
    remainingAmount -= payAmount
    
    // 记录还款详情
    payDetails.append({
        "orderUUID": order.UID,
        "amount": payAmount
    })

// 3. 创建还款记录（包含 payDetails JSON）
```

### 撤销还款算法

```go
// 1. 查询还款记录
pay = SELECT * FROM batch_order_pay WHERE uid = ?

// 2. 根据类型恢复
if pay.BatchOrderUUID != "":
    // 订单还款：直接恢复该订单
    order.CreditAmount += pay.Amount
else if pay.PayDetails != "":
    // 快捷还款：按原还款详情恢复
    for detail in payDetails:
        order.CreditAmount += detail.Amount

// 3. 标记为已撤销
pay.IsRevoked = 1
pay.RevokedAt = now
pay.RevokedReason = req.Reason
```

---

## 📁 文件清单

| 文件 | 说明 |
|------|------|
| `app/service/model/payment.go` | 还款和消息数据模型 |
| `app/service/model/request/payment.go` | 请求参数定义 |
| `app/service/model/response/payment.go` | 响应结构定义 |
| `app/service/logic/payment.go` | 业务逻辑实现 |
| `app/service/handler/payment.go` | HTTP 处理器 |
| `app/service/handler/router.go` | 路由注册 |
| `app/service/model/migrate.go` | 数据库迁移 |

---

## 🚀 使用示例

### 1. 查看快速还款列表

```bash
curl -X POST http://localhost:8080/api/csj/payment/quick/list \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"page": 1, "pageCount": 10}'
```

### 2. 快捷还款 600 元

```bash
curl -X POST http://localhost:8080/api/csj/payment/quick/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customerUUID": "customer-uuid",
    "amount": 600.00,
    "payType": 1,
    "remark": "现金还款"
  }'
```

### 3. 针对订单还款 200 元

```bash
curl -X POST http://localhost:8080/api/csj/payment/order/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "orderUUID": "order-uuid",
    "amount": 200.00,
    "payType": 2,
    "remark": "微信还款"
  }'
```

### 4. 查看还款历史

```bash
curl -X POST http://localhost:8080/api/csj/payment/history \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customerUUID": "customer-uuid",
    "page": 1,
    "pageCount": 10
  }'
```

### 5. 撤销还款

```bash
curl -X POST http://localhost:8080/api/csj/payment/revoke \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "payUUID": "pay-uuid",
    "reason": "还款金额错误"
  }'
```

---

## ⚠️ 注意事项

1. **事务安全**：所有还款操作都使用数据库事务，确保数据一致性
2. **CreditAmount 可以为负**：表示客户多付款了，形成预付款
3. **PayDetails JSON**：快捷还款时记录每个订单的还款明细，方便撤销时恢复
4. **消息中心**：所有还款/撤销操作都会记录到消息中心，便于追溯
5. **幂等性**：已撤销的记录不能再次撤销
6. **权限控制**：所有接口都需要认证，只能操作自己 OwnerUser 下的数据

---

## 🎯 下一步建议

1. 添加还款统计接口（按日/月统计还款金额）
2. 添加消息已读/未读标记功能
3. 添加还款凭证上传功能（图片）
4. 添加还款提醒功能（赊欠超期提醒）
5. 添加导出还款记录功能（Excel）
