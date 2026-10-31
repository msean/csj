# 官方消息与消息红点

## 部署

1. 在 app 与后台共用的 MySQL 数据库执行 `scripts/add_official_messages.sql`（可重复执行）。脚本创建两张表，注册后台菜单和 API，并仅给现有超级管理员角色 `888` 授权。
2. 部署 app、后台和 web，重启后台以重新加载权限，重新登录后进入「官方消息」。其他角色通过现有角色管理分配菜单和 API 权限。
3. 新安装环境也已接入模型迁移和默认菜单/API/权限初始化。app 的 `-m` 模式会创建官方消息及已读表。已有环境仍需执行升级脚本以注册菜单及权限。

本次代码修改没有连接或修改线上数据库。升级前需在目标 MySQL 环境验证脚本及接口。

## 存储与展示

- `official_messages`：官方消息共享一份，发布时不复制到 `message_centers`。
- `official_message_reads`：用户打开官方消息时插入已读记录，联合主键为 `(user_uuid, message_uuid)`，重复调用保持幂等。
- 发布状态 `0` 为草稿，`1` 为发布；发布时间到达、未失效、未删除才对所有用户可见，包括后来注册的用户。
- 所有可见官方消息排在个人消息前，官方消息按 `sort DESC`、发布时间倒序排列；已读后仍置顶。
- 在合并列表上统一分页，不会每页重复插入公告；`page` 默认 1，`pageCount` 默认 20、最大 100。
- 修改原消息内容保留用户已读状态；需要重新提醒时新增一条消息。

## 客户端对接

`POST /message/list` 的原有返回字段保持不变，`items` 新增：

```json
{
  "messageUUID": "官方消息UID",
  "source": "official",
  "title": "更新公告",
  "type": 3,
  "event": "官方消息",
  "content": "公告内容",
  "isRead": 0
}
```

`source` 为 `personal` 或 `official`；官方消息的 `createdAt` 为发布时间。`summary.totalCount` 和 `summary.unreadCount` 包含当前可见官方消息，失效/撤回的官方消息不计入。

用户打开消息后调用 `POST /message/read`：

```json
{"messageUUID":"官方消息UID","source":"official"}
```

不传 `source` 时默认个人消息，兼容旧客户端。官方消息必须传 `source: "official"`。不可见或不存在的官方消息不能被标记已读。获取列表本身不会标记已读。

用户 profile 新增 `hasUnreadMessage: boolean`，判断个人未读消息或当前可见官方未读消息是否存在。前端用它显示红点，在已读成功和页面恢复时刷新；此字段不提供实时推送。

## 验证

- `cd app && go test ./service ./service/dao ./service/logic ./service/handler ./service/model/...`
- `cd backend && go test ./model/appmgr/... ./service/appmgr ./api/appmgr ./router/appmgr ./initialize ./source/system`
- `python3 -B -m unittest discover -s scripts/tests -v`：运行生产列表 SQL 和可见性条件，验证置顶、跨源分页、用户隔离、草稿/未来/失效/删除过滤及已读状态。此测试使用 SQLite，不替代目标 MySQL 的集成验证。
- `cd web && npm run build`
