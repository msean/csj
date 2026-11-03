# 官方通知与客服配置

部署前在 app/backend 共用的数据库执行 `scripts/notification_support.sql`（可重复执行）。新增通知等级字段、客服配置表和后台客服 GET/PUT API 权限。默认给已有超级管理员角色 888 授权，其他角色按需在 API 权限中授权。已有数据库须运行该脚本以登记 API；模型 AutoMigrate 不能代替权限初始化。执行后重启 backend，令权限缓存重新加载，再部署 app 服务、web 和 clerk。

后台「App业务管理 → 客服配置」可直接维护客服信息，官方消息页面顶部也保留快捷入口。

后台「官方消息」：新增/编辑可选择一般（level=1，橙色）或重要（level=2，红色），旧消息默认为一般。顶部「联系客服配置」维护客服电话、客服微信及说明/工作时间；留空保存可清除。未配置时移动端显示暂未配置。

移动端首页喇叭显示排序值最大、发布时间最新的有效通知；点击进入独立分页列表。无有效通知时隐藏喇叭。官方通知不展示或写入已读状态，列表与详情展示发布时间，标题和正文预览各最多两行。通知仍为共享存储；草稿、未到发布时间、过期、删除的通知不展示。点击列表项直接查看完整内容，不再等待已读请求。

个人消息中心及 profile.hasUnreadMessage 只统计个人业务消息。POST /api/csj/message/readAll 按当前登录用户批量标记个人消息已读，完成后移动端刷新列表并同步铃铛。铃铛在有未读时为黄色、摇摆并展示较大红点，无未读时恢复。官方通知不进入个人消息中心。

移动接口：
- POST /api/csj/message/official/list：登录后分页，参数 page/pageCount，返回 items 与 summary.totalCount/unreadCount。
- POST /api/csj/message/official/read：仅为旧版兼容保留，新版不调用。
- POST /api/csj/support/config：只读公开客服联系方式，无需登录。

后台接口：GET/PUT /officialMessage/supportConfig，均使用已有后台鉴权、API 权限及 PUT 操作日志。

clerk 是独立 Git 仓库，移动端修改需在 clerk 内单独检查、提交及打包。
