-- Existing installations: run add_official_messages.sql first, then this script.
-- Does not change business records or create order shards.
-- API permissions are granted only to the existing super-admin role (888).
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 0, 0, 'appManagement', 'appManagement', 0, 'view/routerHolder.vue', 10, 'App业务管理', 'management'
WHERE NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE name = 'appManagement');
SET @app_management_menu_id = (SELECT id FROM sys_base_menus WHERE name = 'appManagement' AND deleted_at IS NULL LIMIT 1);

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'appUser', 'AppUser', 0, 'view/appUser/index.vue', 1, '注册用户', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appUser/index.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appUser/index.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'appVersion', 'AppVersion', 0, 'view/appVersion/index.vue', 2, 'App版本', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appVersion/index.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appVersion/index.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'memberConfig', 'MemberConfig', 0, 'view/memberConfig/index.vue', 3, '会员配置', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/memberConfig/index.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/memberConfig/index.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'rechargeOrder', 'RechargeOrder', 0, 'view/rechargeOrder/index.vue', 4, '充值订单', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/rechargeOrder/index.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/rechargeOrder/index.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'officialMessage', 'officialMessage', 0, 'view/officialMessage/index.vue', 5, '官方消息', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/officialMessage/index.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/officialMessage/index.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'customers', 'AppCustomers', 0, 'view/appData/AppCustomers.vue', 6, '客户管理', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppCustomers.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppCustomers.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'goods-categories', 'AppGoodsCategories', 0, 'view/appData/AppGoodsCategories.vue', 7, '商品分类', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppGoodsCategories.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppGoodsCategories.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'goods', 'AppGoods', 0, 'view/appData/AppGoods.vue', 8, '商品管理', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppGoods.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppGoods.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'batches', 'AppBatches', 0, 'view/appData/AppBatches.vue', 9, '批次管理', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppBatches.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppBatches.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'orders', 'AppOrders', 0, 'view/appData/AppOrders.vue', 10, '订单查询', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppOrders.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppOrders.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'payments', 'AppPayments', 0, 'view/appData/AppPayments.vue', 11, '还款查询', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppPayments.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppPayments.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'messages', 'AppMessages', 0, 'view/appData/AppMessages.vue', 12, '个人消息', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppMessages.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppMessages.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_management_menu_id, 'shares', 'AppShares', 0, 'view/appData/AppShares.vue', 13, '分享记录', 'tickets'
WHERE @app_management_menu_id IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE component = 'view/appData/AppShares.vue' AND deleted_at IS NULL);
UPDATE sys_base_menus SET parent_id = @app_management_menu_id, menu_level = 1
WHERE component = 'view/appData/AppShares.vue' AND deleted_at IS NULL AND @app_management_menu_id IS NOT NULL;

INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT m.id, 888 FROM sys_base_menus m
WHERE (m.id = @app_management_menu_id OR m.parent_id = @app_management_menu_id)
AND m.deleted_at IS NULL AND EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM sys_authority_menus a WHERE a.sys_base_menu_id = m.id AND a.sys_authority_authority_id = 888);

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appUser/findAppUser', '注册用户', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appUser/findAppUser' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appUser/findAppUser', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appUser/findAppUser' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appUser/getAppUserList', '注册用户', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appUser/getAppUserList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appUser/getAppUserList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appUser/getAppUserList' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/createAppVersion', 'App版本', 'App业务管理', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/createAppVersion' AND method = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/createAppVersion', 'POST', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/createAppVersion' AND v2 = 'POST');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/deleteAppVersion', 'App版本', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/deleteAppVersion' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/deleteAppVersion', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/deleteAppVersion' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/deleteAppVersionByIds', 'App版本', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/deleteAppVersionByIds' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/deleteAppVersionByIds', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/deleteAppVersionByIds' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/updateAppVersion', 'App版本', 'App业务管理', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/updateAppVersion' AND method = 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/updateAppVersion', 'PUT', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/updateAppVersion' AND v2 = 'PUT');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/findAppVersion', 'App版本', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/findAppVersion' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/findAppVersion', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/findAppVersion' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appVersion/getAppVersionList', 'App版本', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appVersion/getAppVersionList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appVersion/getAppVersionList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appVersion/getAppVersionList' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/createMemberConfig', '会员配置', 'App业务管理', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/createMemberConfig' AND method = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/createMemberConfig', 'POST', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/createMemberConfig' AND v2 = 'POST');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/deleteMemberConfig', '会员配置', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/deleteMemberConfig' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/deleteMemberConfig', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/deleteMemberConfig' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/deleteMemberConfigByIds', '会员配置', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/deleteMemberConfigByIds' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/deleteMemberConfigByIds', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/deleteMemberConfigByIds' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/updateMemberConfig', '会员配置', 'App业务管理', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/updateMemberConfig' AND method = 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/updateMemberConfig', 'PUT', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/updateMemberConfig' AND v2 = 'PUT');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/findMemberConfig', '会员配置', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/findMemberConfig' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/findMemberConfig', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/findMemberConfig' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/memberConfig/getMemberConfigList', '会员配置', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/memberConfig/getMemberConfigList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/memberConfig/getMemberConfigList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/memberConfig/getMemberConfigList' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/rechargeOrder/findRechargeOrder', '充值订单', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/rechargeOrder/findRechargeOrder' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/rechargeOrder/findRechargeOrder', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/rechargeOrder/findRechargeOrder' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/rechargeOrder/getRechargeOrderList', '充值订单', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/rechargeOrder/getRechargeOrderList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/rechargeOrder/getRechargeOrderList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/rechargeOrder/getRechargeOrderList' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/createOfficialMessage', '官方消息', 'App业务管理', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/createOfficialMessage' AND method = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/createOfficialMessage', 'POST', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/createOfficialMessage' AND v2 = 'POST');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/deleteOfficialMessage', '官方消息', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/deleteOfficialMessage' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/deleteOfficialMessage', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/deleteOfficialMessage' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/deleteOfficialMessageByIds', '官方消息', 'App业务管理', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/deleteOfficialMessageByIds' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/deleteOfficialMessageByIds', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/deleteOfficialMessageByIds' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/updateOfficialMessage', '官方消息', 'App业务管理', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/updateOfficialMessage' AND method = 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/updateOfficialMessage', 'PUT', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/updateOfficialMessage' AND v2 = 'PUT');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/findOfficialMessage', '官方消息', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/findOfficialMessage' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/findOfficialMessage', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/findOfficialMessage' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/getOfficialMessageList', '官方消息', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/getOfficialMessageList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/getOfficialMessageList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/getOfficialMessageList' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/customers/list', '客户管理列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/customers/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/customers/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/customers/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/customers/detail', '客户管理详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/customers/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/customers/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/customers/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/goods-categories/list', '商品分类列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/goods-categories/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/goods-categories/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/goods-categories/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/goods-categories/detail', '商品分类详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/goods-categories/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/goods-categories/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/goods-categories/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/goods/list', '商品管理列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/goods/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/goods/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/goods/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/goods/detail', '商品管理详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/goods/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/goods/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/goods/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/batches/list', '批次管理列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/batches/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/batches/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/batches/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/batches/detail', '批次管理详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/batches/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/batches/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/batches/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/orders/list', '订单查询列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/orders/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/orders/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/orders/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/orders/detail', '订单查询详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/orders/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/orders/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/orders/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/payments/list', '还款查询列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/payments/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/payments/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/payments/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/payments/detail', '还款查询详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/payments/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/payments/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/payments/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/messages/list', '个人消息列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/messages/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/messages/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/messages/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/messages/detail', '个人消息详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/messages/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/messages/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/messages/detail' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/shares/list', '分享记录列表', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/shares/list' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/shares/list', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/shares/list' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/appData/shares/detail', '分享记录详情', 'App业务管理', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/appData/shares/detail' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/appData/shares/detail', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/appData/shares/detail' AND v2 = 'GET');
