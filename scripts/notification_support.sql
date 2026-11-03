-- Run against the shared app/backend database before deploying the new binaries.
-- Safe to rerun. Existing notifications default to general (orange).
SET @has_level = (SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'official_messages' AND column_name = 'level');
SET @ddl = IF(@has_level = 0, 'ALTER TABLE official_messages ADD COLUMN level bigint NOT NULL DEFAULT 1', 'SELECT 1');
PREPARE notification_stmt FROM @ddl;
EXECUTE notification_stmt;
DEALLOCATE PREPARE notification_stmt;
CREATE TABLE IF NOT EXISTS support_configs (
 id bigint unsigned NOT NULL PRIMARY KEY,
 phone varchar(50) NOT NULL DEFAULT '',
 wechat varchar(100) NOT NULL DEFAULT '',
 description varchar(2000) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/supportConfig', '客服配置查询', '官方消息', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/officialMessage/supportConfig' AND method='GET' AND deleted_at IS NULL);
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p','888','/officialMessage/supportConfig','GET','','',''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id=888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0='888' AND v1='/officialMessage/supportConfig' AND v2='GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/supportConfig', '客服配置保存', '官方消息', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path='/officialMessage/supportConfig' AND method='PUT' AND deleted_at IS NULL);
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p','888','/officialMessage/supportConfig','PUT','','',''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id=888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype='p' AND v0='888' AND v1='/officialMessage/supportConfig' AND v2='PUT');

-- Dedicated customer support configuration page under the existing App menu.
SET @app_menu = (SELECT id FROM sys_base_menus WHERE name='appManagement' AND deleted_at IS NULL LIMIT 1);
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1, @app_menu, 'supportConfig', 'SupportConfig', 0, 'view/supportConfig/index.vue', 14, '客服配置', 'service'
WHERE @app_menu IS NOT NULL
AND NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE (name='SupportConfig' OR component='view/supportConfig/index.vue') AND deleted_at IS NULL);
INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT m.id, 888 FROM sys_base_menus m
WHERE m.component='view/supportConfig/index.vue' AND m.deleted_at IS NULL
AND EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id=888)
AND NOT EXISTS (SELECT 1 FROM sys_authority_menus a WHERE a.sys_base_menu_id=m.id AND a.sys_authority_authority_id=888);
