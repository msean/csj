-- ======================================================================
-- 统一当前库所有表的 collation 为 utf8mb4_unicode_ci
--
-- 两步走：
--   1) ALTER TABLE ... COLLATE utf8mb4_unicode_ci   —— 改表的 default collation（影响后续新增列）
--   2) ALTER TABLE ... MODIFY COLUMN ...            —— 把已有字符串列的 collation 改成 utf8mb4_unicode_ci
--
-- 注意：
--   1. utf8mb4 → utf8mb4 同 charset，只改 collation 不会丢数据
--   2. 大表 MODIFY COLUMN 会重建索引，耗时较长，建议低峰期执行
--   3. 执行前建议备份：mysqldump -u root -p caishuji > backup.sql
-- ======================================================================

-- 放宽 group_concat_max_len
SET SESSION group_concat_max_len = 1048576;

DROP PROCEDURE IF EXISTS fix_all_collation;
DELIMITER $$
CREATE PROCEDURE fix_all_collation()
fix_all: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl VARCHAR(64);
    DECLARE col_sql LONGTEXT;
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_TYPE = 'BASE TABLE'
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    OPEN cur;
    tbl_loop: LOOP
        FETCH cur INTO tbl;
        IF done = 1 THEN LEAVE tbl_loop; END IF;

        -- 1) 改表的 default collation
        SET @ddl := CONCAT('ALTER TABLE `', tbl, '` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci');
        PREPARE stmt FROM @ddl;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;

        -- 2) 收集该表所有非 utf8mb4_unicode_ci 的字符串列，拼成 MODIFY COLUMN 语句
        --    保留原数据类型 / nullable / 默认值 / 注释
        --    information_schema.COLUMNS 已经把类型拆成 DATA_TYPE / CHARACTER_MAXIMUM_LENGTH 等
        --    我们用 COLUMN_TYPE 直接拿到完整类型定义（如 varchar(64) / longtext）
        SELECT GROUP_CONCAT(
            CONCAT(
                'MODIFY COLUMN `', COLUMN_NAME, '` ',
                COLUMN_TYPE,
                ' CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci ',
                IF(IS_NULLABLE = 'NO', 'NOT NULL ', 'NULL '),
                IF(COLUMN_DEFAULT IS NULL, '', CONCAT('DEFAULT ', QUOTE(COLUMN_DEFAULT), ' ')),
                IF(COLUMN_COMMENT = '', '', CONCAT('COMMENT ', QUOTE(COLUMN_COMMENT), ' '))
            )
            SEPARATOR ', '
        ) INTO col_sql
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME = tbl
          AND DATA_TYPE IN ('varchar', 'char', 'text', 'longtext', 'mediumtext', 'tinytext', 'long', 'medium', 'tiny')
          AND CHARACTER_SET_NAME IS NOT NULL
          AND COLLATION_NAME <> 'utf8mb4_unicode_ci';

        IF col_sql IS NOT NULL AND col_sql != '' THEN
            SET @ddl := CONCAT('ALTER TABLE `', tbl, '` ', col_sql);
            PREPARE stmt FROM @ddl;
            EXECUTE stmt;
            DEALLOCATE PREPARE stmt;
        END IF;

    END LOOP;
    CLOSE cur;

    SELECT '全部表 collation 统一完成' AS msg;
END$$
DELIMITER ;

-- 执行
CALL fix_all_collation();

-- 清理
DROP PROCEDURE IF EXISTS fix_all_collation;

-- ==================== 验证 ====================
-- 1) 表的 default collation（应该全部是 utf8mb4_unicode_ci）
SELECT TABLE_NAME, TABLE_COLLATION
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_TYPE = 'BASE TABLE'
  AND TABLE_COLLATION <> 'utf8mb4_unicode_ci'
ORDER BY TABLE_NAME;
-- 上面查询应该返回 0 行

-- 2) 列的 collation（应该全部是 utf8mb4_unicode_ci）
SELECT TABLE_NAME, COLUMN_NAME, COLLATION_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND CHARACTER_SET_NAME IS NOT NULL
  AND COLLATION_NAME <> 'utf8mb4_unicode_ci'
ORDER BY TABLE_NAME, COLUMN_NAME;
-- 上面查询应该返回 0 行
