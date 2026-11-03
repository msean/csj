# 本次移动接口更新

移动端使用的是 app 服务（/api/csj），不是后台管理 backend 服务。

已在 193.112.131.40 的 HTTP 入口只读核验：POST /api/csj/batch_order/customer_purchases 和 POST /api/csj/message/readAll 返回 404，已有消息列表接口能进入鉴权。这表明当前入口尚未加载新路由。仅更新 clerk、后台 backend，或重启仍使用旧镜像的容器，不能让新接口生效。需要构建并部署包含这些修改的 app 服务，随后更新 clerk。本次没有替换线上服务。

采购汇总已修正为实际数据库列 amount（模型对外字段名 totalAmount），并在真实 MySQL 分表上用不存在的测试标识只读执行验证。无需修改数据库表结构。

POST /api/csj/payment/detail 原响应新增 allocations 数组；原 payDetails 保留兼容。每项包含 type、orderUUID、amount（本笔抵扣金额）、orderAvailable；关联订单可用时另包含 orderTime（开单时间，老订单回退创建时间）、orderTotalAmount（订单金额加三轮车费）、orderStatus。关联数据按记录所属账号及客户限定，批量查询；历史欠款、已删除或无法匹配的订单不补造金额和时间。

这些订单信息来自关联订单的当前记录，不是还款当时的快照。接口只读，不修改原有还款分配和历史账务数据。
