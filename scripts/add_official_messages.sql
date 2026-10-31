-- Run on the shared app/admin MySQL database before deploying this feature.
CREATE TABLE IF NOT EXISTS official_messages (
 uid varchar(64) NOT NULL PRIMARY KEY,
 created_at datetime(3) NULL,
 updated_at datetime(3) NULL,
 deleted_at datetime(3) NULL,
 title varchar(200) NOT NULL,
 content text NOT NULL,
 status bigint NOT NULL DEFAULT 0,
 published_at datetime(3) NOT NULL,
 expires_at datetime(3) NULL,
 sort bigint NOT NULL DEFAULT 0,
 KEY idx_official_messages_status (status),
 KEY idx_official_messages_published_at (published_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
CREATE TABLE IF NOT EXISTS official_message_reads (
 user_uuid varchar(64) NOT NULL,
 message_uuid varchar(64) NOT NULL,
 read_at datetime(3) NOT NULL,
 PRIMARY KEY (user_uuid, message_uuid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Register the menu and APIs; only the existing super-admin role (888) is granted access.
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 0, 0, 'officialMessage', 'officialMessage', 0, 'view/officialMessage/index.vue', 10, '官方消息', 'message'
WHERE NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE name = 'officialMessage');
INSERT INTO sys_authority_menus (sys_base_menu_id, sys_authority_authority_id)
SELECT m.id, 888 FROM sys_base_menus m
WHERE m.name = 'officialMessage' AND m.deleted_at IS NULL
AND EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM sys_authority_menus a WHERE a.sys_base_menu_id = m.id AND a.sys_authority_authority_id = 888);

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/createOfficialMessage', '官方消息创建', '官方消息', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/createOfficialMessage' AND method = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/createOfficialMessage', 'POST', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/createOfficialMessage' AND v2 = 'POST');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/updateOfficialMessage', '官方消息更新', '官方消息', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/updateOfficialMessage' AND method = 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/updateOfficialMessage', 'PUT', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/updateOfficialMessage' AND v2 = 'PUT');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/deleteOfficialMessage', '官方消息删除', '官方消息', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/deleteOfficialMessage' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/deleteOfficialMessage', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/deleteOfficialMessage' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/deleteOfficialMessageByIds', '官方消息批量删除', '官方消息', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/deleteOfficialMessageByIds' AND method = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/deleteOfficialMessageByIds', 'DELETE', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/deleteOfficialMessageByIds' AND v2 = 'DELETE');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/findOfficialMessage', '官方消息详情', '官方消息', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/findOfficialMessage' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/findOfficialMessage', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/findOfficialMessage' AND v2 = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/officialMessage/getOfficialMessageList', '官方消息列表', '官方消息', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/officialMessage/getOfficialMessageList' AND method = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3, v4, v5)
SELECT 'p', '888', '/officialMessage/getOfficialMessageList', 'GET', '', '', ''
WHERE EXISTS (SELECT 1 FROM sys_authorities WHERE authority_id = 888)
AND NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = '888' AND v1 = '/officialMessage/getOfficialMessageList' AND v2 = 'GET');
