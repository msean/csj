# 会员充值微信支付接入实现计划

## Repository Research

### 项目现状
- **app（移动端后端）**：Gin 框架，使用 GORM + MySQL，JWT 认证。现有 `User` 模型含 `ExperienceDeadline`（体验截止时间，列名 `vip_expire_time`）。现有支付逻辑仅为业务层面的还款记录（`BatchOrderPay`），未接入真实微信支付。
- **backend（管理后台）**：Gin-Vue-Admin 框架，复用 app 的表结构（uid 主键）。
- **web（管理前端）**：Vue3 + Element Plus。
- **会员配置**：`member_configs` 表已创建，type（1月/2季度/3年）+ price。

### 用户确认的决策
1. **支付方式**：微信 App 支付（原生 App 内调用微信客户端）
2. **会员截止时间**：在 `users` 表新增独立字段 `vip_deadline`，与体验期分开
3. **展示端**：web 管理后台展示所有用户充值记录和会员截止时间
4. **订单用途**：充值订单仅用于会员开通/续费

### 技术选型
- 微信支付 SDK：`github.com/wechatpay-apiv3/wechatpay-go`（官方 V3 SDK）
- 微信 App 支付流程：后端统一下单 → 返回签名参数给 App → App 调起微信支付 → 微信异步回调后端 → 后端验签并延长会员时间

## Files and Modules

### 数据库变更
- `app/service/model/user.go`：User 模型新增 `VipDeadline` 字段
- `app/service/model/recharge_order.go`（新增）：充值订单模型
- `app/service/model/migrate.go`：注册新模型，AutoMigrate 自动加字段
- `backend/model/appmgr/recharge_order.go`（新增）：backend 侧充值订单模型（uid 主键）
- `backend/model/appmgr/user.go` 或已有用户模型：新增 vip_deadline 字段（如 backend 有 User 模型）

### app 后端
- `deploy/app_debug_conf.yaml`：新增 wechat 配置节
- `app/pkg/wechat/pay.go`（新增）：微信支付服务封装（统一下单、回调验签、订单查询）
- `app/service/model/request/recharge.go`（新增）：请求参数
- `app/service/model/response/recharge.go`（新增）：响应结构
- `app/service/dao/recharge.go`（新增）：充值订单 DAO
- `app/service/logic/recharge.go`（新增）：充值业务逻辑
- `app/service/handler/recharge.go`（新增）：充值路由和 handler
- `app/service/handler/router.go`：注册充值路由
- `app/service/logic/user.go` / `app/service/model/response/user.go`：用户信息返回会员截止时间

### backend 后端
- `backend/model/appmgr/recharge_order.go`（新增）
- `backend/model/appmgr/request/recharge_order.go`（新增）
- `backend/service/appmgr/recharge_order.go`（新增）
- `backend/api/appmgr/recharge_order.go`（新增）
- `backend/router/appmgr/recharge_order.go`（新增）
- `backend/api/appmgr/enter.go`、`service/appmgr/enter.go`、`router/appmgr/enter.go`：注册
- `backend/initialize/ensure_tables.go`：注册 RechargeOrder 表
- `backend/initialize/router.go`：注册路由

### web 前端
- `web/src/api/rechargeOrder.js`（新增）：充值订单 API
- `web/src/view/rechargeOrder/index.vue`（新增）：充值记录管理页面
- 用户管理页面：展示会员截止时间（如已有用户管理页面）

## Implementation Steps

### 阶段一：数据模型与配置
1. 在 `app/service/model/user.go` 的 User 结构体新增 `VipDeadline time.Time` 字段，列名 `vip_deadline`
2. 新建 `app/service/model/recharge_order.go`：RechargeOrder 模型（uid、user_uuid、order_no 唯一、member_type、amount、pay_status、pay_time、transaction_id、时间戳）
3. 在 `app/service/model/migrate.go` 注册 `&RechargeOrder{}`
4. 在 `deploy/app_debug_conf.yaml` 新增 wechat 配置节（appid、mch_id、api_v3_key、cert_serial_no、private_key_path、notify_url）

### 阶段二：微信支付服务封装
5. 执行 `go get github.com/wechatpay-apiv3/wechatpay-go` 添加依赖
6. 新建 `app/pkg/wechat/pay.go`：
   - 初始化微信支付 Client（读取配置）
   - `AppPay`：统一下单，返回 App 支付所需的签名参数
   - `VerifyNotify`：回调验签
   - `QueryOrder`：查询订单状态

### 阶段三：app 充值业务逻辑
7. 新建 `app/service/model/request/recharge.go` 和 `response/recharge.go`
8. 新建 `app/service/dao/recharge.go`：订单创建、按订单号查询、更新支付状态、按用户查询列表
9. 新建 `app/service/logic/recharge.go`：
   - `CreateRechargeOrder`：选择会员类型 → 查 member_configs 价格 → 创建订单（order_no 唯一）→ 调用微信 App 支付统一下单 → 返回支付参数
   - `WechatNotify`：验签 → 查订单 → 幂等校验（已支付则直接返回成功）→ 更新订单状态（pay_status=1, pay_time, transaction_id）→ 延长用户 vip_deadline（按 member_type 加 1月/3月/1年）
   - `QueryRechargeStatus`：查询订单支付状态（供前端轮询）
   - `RechargeList`：当前用户充值记录列表
10. 新建 `app/service/handler/recharge.go`：
    - `POST /api/csj/recharge/create`：创建充值订单
    - `POST /api/csj/recharge/notify`：微信支付回调（无需鉴权）
    - `POST /api/csj/recharge/status`：查询订单支付状态
    - `POST /api/csj/recharge/list`：充值记录列表
11. 在 `app/service/handler/router.go` 注册充值路由
12. 修改 `app/service/logic/user.go` 的 Profile 方法，返回 vip_deadline

### 阶段四：backend 管理端充值记录
13. 新建 `backend/model/appmgr/recharge_order.go`（uid 主键，与 app 表结构一致）
14. 新建 request/service/api/router（CRUD，主要是列表查询和详情）
15. 在 enter.go、ensure_tables.go、router.go 注册

### 阶段五：web 前端页面
16. 新建 `web/src/api/rechargeOrder.js`
17. 新建 `web/src/view/rechargeOrder/index.vue`：充值记录列表（可按用户、状态、时间筛选），展示订单号、用户、会员类型、金额、支付状态、支付时间、交易号

### 阶段六：编译验证
18. `go build ./...` 验证 app 和 backend 编译通过
19. `npm run build` 验证 web 编译通过

## Dependencies and Considerations

- **微信支付配置**：商户号、APIv3 密钥、商户证书等敏感信息通过配置文件读取，不硬编码。部署时需用户在 yaml 中填写。
- **会员时间延长逻辑**：若用户当前 vip_deadline 在未来，则从 vip_deadline 起延长；若已过期或为空，则从当前时间起延长。
- **幂等性**：微信回调可能重复通知，必须以 order_no 为维度做幂等校验，已支付订单不再重复延长会员时间。
- **回调验签**：必须使用微信平台证书验签回调请求，确保回调来自微信。
- **订单号生成**：使用时间戳+随机数生成，确保全局唯一。
- **backend 与 app 共享表**：backend 的 RechargeOrder 模型使用 uid 主键，与 app 一致（参照 AppVersion 的处理方式）。

## Validation

- app `go build ./...` 编译通过
- backend `go build ./...` 编译通过（受影响包）
- web `npm run build` 编译通过
- 数据库迁移成功创建 recharge_orders 表，users 表新增 vip_deadline 字段

## Risks

- **微信支付 SDK 依赖**：需联网下载 `wechatpay-go` 依赖，若网络受限需手动 vendor。
- **微信支付证书**：需要用户提供商户证书文件（apiclient_key.pem），代码中只读取路径，不内置证书。
- **回调地址公网可达**：微信支付回调要求 notify_url 为公网可访问地址，本地开发需使用内网穿透。
- **backend 用户模型**：需确认 backend 是否有 User 模型，若有则需同步添加 vip_deadline 字段用于展示。
