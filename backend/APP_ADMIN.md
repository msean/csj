# App 业务后台

## 页面范围

后台「App业务管理」包含 13 个入口：注册用户、App 版本、会员配置、充值订单、官方消息、客户、商品分类、商品、批次、订单、还款、个人消息、分享记录。

- 版本、会员配置、官方消息保留创建、编辑和软删除功能。
- 用户、客户、分类、商品、批次、订单、还款、个人消息、充值订单、分享记录提供查询和详情；业务数据的写入继续使用 app 的业务接口，保持账务、库存、操作历史与缓存更新的一致性。
- 用户详情可以跳转到该用户名下各类数据；客户详情可以跳转到客户订单、还款；批次详情可以跳转到该批次订单。
- 批次详情包括批次商品；订单详情包括订单商品和操作历史；还款详情包括客户信息及完整还款字段。
- 不提供验证码、计数器等内部基础设施表的后台编辑入口。

## 查询接口

所有新接口位于后台原有 JWT 和 Casbin 权限组下：

```
GET /appData/{resource}/list
GET /appData/{resource}/detail?uid=...&ownerUser=...
```

`resource` 只接受 `customers`、`goods-categories`、`goods`、`batches`、`orders`、`payments`、`messages`、`shares`。

列表支持 `page/pageSize/ownerUser/keyword/from/to`；支持状态的页面使用可选 `status`，状态 `0` 不会当成“未筛选”。订单、还款、消息、分享按各自模型支持 `customerUUID`；订单也支持 `batchUUID`。时间使用 RFC3339，页码默认 1，页大小默认 10、最大 100。

订单及还款必须选择所属用户，用共享分表函数定位对应表，并额外按 `owner_user` 限定；不跨 50 张表扫描。详情要求所属用户，子记录也限定到已验证的记录。接口只接受资源白名单，不接受用户提供表名。

原接口 JSON 命名与路径保持兼容。充值查询修复了 `payStatus=0` 待支付筛选，原有列表统一了分页默认值和上限。

## 部署

1. 已有数据库依次执行 `scripts/add_official_messages.sql` 和 `scripts/complete_app_admin.sql`。脚本可重复运行，后者注册/归组菜单、API 和超级管理员 888 的权限，不修改业务记录、不自动创建订单分表。
2. 部署 app、backend 和 web；重启 backend 重新加载权限，重新登录后台。
3. 其他角色通过原角色管理分配相应菜单和 API。若使用所属用户搜索框，也需要 `/appUser/getAppUserList` 查询权限。

新建数据库已接入默认菜单/API/Casbin 初始化。app 的业务表与分表仍按原有流程迁移。

本次开发没有对真实数据库运行升级脚本；目标 MySQL 环境仍需进行部署验证。

## 开发验证

```bash
python3 scripts/sync_shared_vendor.py --check
(cd app && go test ./service/... ./pkg/utils)
(cd backend && go test ./model/appmgr/... ./service/appmgr ./api/appmgr ./router/appmgr ./initialize ./source/system)
python3 -B -m unittest discover -s scripts/tests -v
```

新增测试覆盖用户隔离、关键词 OR 条件隔离、分页、资源白名单、分表查询、关联明细隔离、待支付筛选、共享模型删除行为及 JSON 兼容性。SQL 行为测试使用 SQLite，不能替代目标 MySQL 验证。

前端验证使用 `vite build --minify esbuild --outDir /private/tmp/caishuji-admin-web-build`，构建兼容与现代浏览器两组产物；同时编译检查所有相关 Vue 组件。构建有 CSS 压缩和大包体积告警，未改项目默认的 terser 构建配置。

校验期间修复了原 JWT 测试依赖未初始化运行时/固定过期令牌的问题，以及原公告代码生成器 `go:generate` 指令放在 import 行尾的问题；未修改 JWT 的生产逻辑。
