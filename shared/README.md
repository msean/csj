# 共享业务模型

`github.com/msean/csj/shared` 是仓库内的独立 Go 模块，通过两个服务各自 `go.mod` 中的 `require` / `replace ../shared` 引用。无需把 backend 依赖到 app 的启动代码、配置、数据库连接、缓存或 HTTP 服务，也无需发布到第三方平台。

- `model/`：唯一的业务实体定义，保留 app 的表结构、字段标签、ULID 主键、GORM 钩子和业务方法。实体显式声明原有表名，避免后台的表前缀设置改变 app 表的定位。
- `constants/`：业务枚举定义。
- `sharding/`：统一 CRC32 / 50 分表路由，backend 不自行猜测表名。
- `internal/modelutil/`：模型方法依赖的数值格式化函数。
- app 的 `service/model/entities.go` 与 backend 的 `model/appmgr/` 仅保留类型别名，不再维护第二套持久化结构。
- app 的请求/响应 DTO、迁移执行、数据库连接和服务逻辑仍归 app 所有；后台的请求/响应转换归 backend 所有。

app 的 JSON 使用 `uuid/createTime/updateTime/vipRemainDays`，后台继续使用 `uid/createdAt/updatedAt/experienceDeadline`。后台通过 `api/appmgr/record_codec.go` 进行边界转换，实体不因管理端命名而改变。

`BaseModel.DeletedAt` 保留 app 原有的 `*time.Time` 语义。后台配置服务显式更新删除时间并过滤已删除记录，避免将原后台软删除意外变成物理删除。app 的会员配置和版本查询也过滤已删除配置。

## 修改与构建

本仓库原本提交了 app/backend 两套 vendor，因此仍支持离线的 vendor 构建。vendor 中的共享包是生成快照，不是另一份编辑入口。

```bash
# 在仓库根目录修改 shared/ 后执行
python3 scripts/sync_shared_vendor.py
python3 scripts/sync_shared_vendor.py --check

# 两个服务独立构建，不需要 go.work
(cd app && go build -mod=vendor ./...)
(cd backend && go build -mod=vendor ./...)
```

同步脚本仅维护共享包及 backend 新增的 ULID 依赖，不改其他 vendor 依赖。将来更新外部依赖时，在正常联网环境执行各服务的 `go mod tidy` / `go mod vendor`，再验证 `--check`。backend Dockerfile 使用提交的 vendor，构建上下文仍可只包含 backend 目录。

## 验证

`backend/model/appmgr/shared_contract_test.go` 验证共享类型、ULID、默认订单状态、金额舍入与分表路由；backend API 测试验证两端 JSON 命名兼容；app/backend 均进行编译检查。数据库结构迁移仍由 app 的迁移流程管理，不由后台数据浏览接口创建或修改业务表。
