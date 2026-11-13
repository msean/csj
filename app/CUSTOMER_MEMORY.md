# 客户快捷开单

`POST /api/csj/batch_order/customer_memory`，沿用现有登录鉴权。

请求：`{"customerUUID":"客户UUID"}`。

响应 `data` 为 `{"orderUUID":"上一单UUID","goodsList":[]}`。客户必须属于当前登录商户且不是现金客户。仅查询未删除、状态为 1/2/3 的订单，按 `order_time DESC, uid DESC` 取最近一单；无历史订单返回空数组与空 orderUUID。

货品字段沿用订单明细（goodsUUID、name、type、goodsWeight、price 等），mount、weight（本次总销售重量）、total 都为 0。goodsWeight 为包装重量，不归零。接口只查询，不创建订单。

clerk 从快捷开单选客户后直接进入购物车，保留上一单单价，数量全部归零；费用、下单时间不沿用。没有历史订单时进入空购物车。快捷开单可在已选货品时进入；进入前保存普通开单草稿，返回时恢复原货品、数量、价格、客户、费用和时间，快捷订单货品不带回。普通选货开单默认现金客户；购物车更换买家时保留当前货品，不恢复目标客户草稿。购物车已有零数量校验，必须填写数量或删除该行后才允许提交。

## 开单使用精简批次详情

`POST /api/csj/batch/detail` 请求可传 `{"uuid":"批次UUID","skipStats":true}`。

- `skipStats` 省略或为 `false`：保留原有批次和货品统计计算及字段，供批次详情和统计页面使用。
- `skipStats: true`：加载批次、货品及名称/类型/包装规格，并单独汇总有效订单的件数和重量计算库存；跳过 `BatchGoodsStat` 的销售金额、均价和利润统计。响应不含批次的 `stat*` 字段，返回货品 `surplus` 库存字段，不含 `sellAmount`、`sellAvgPrice`、`sellTotalAmount`、`sellProfit`。
- 不传 uuid 时查询最新批次，也遵循此参数。

clerk 的 DealController 使用精简模式，BatchController、SalesDetailController、PurchaseOrderDetailController 保持完整模式。精简模式返回剩余库存，开单卡片在绿底规格/价格标签下一行显示库存。定装库存单位为件，散装为斤；未售与负库存均保留。

## 订单货品顺序

订单货品表及其全部分表新增 `sort_no INT NOT NULL DEFAULT 0`（JSON 字段 `sortNo`）。记账单、下单和改单保存时，服务端按照前端 `goodsList` 数组从 1 开始重新编号，不信任客户端传入的序号。重复货品明细也分别保留位置。

快捷开单、订单详情、订单列表及分享单的货品统一按 `sort_no ASC, created_at ASC, uid ASC` 读取。前端无需新增参数。原有记录默认 0，以创建时间与明细 ID 提供稳定顺序；无法还原过去未保存的人工顺序。新一次开单保存后延续本次顺序，新增货品按数组位置保存。

部署前必须为已有订单货品分表补齐字段：使用现有 `--shard-migrate`（`-s`）启动参数及原配置。迁移会检查各分表，只为缺少字段的已有分表添加 `sort_no`，可重复执行；不会回写或重排历史明细。新建分表 SQL 和共享模型也已包含字段。未执行迁移前不要让新版服务接收请求，否则查询排序和写入会缺少列。
