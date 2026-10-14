-- ======================================================================
-- 统一分表的字符串列排序规则为 utf8mb4_unicode_ci，与 customers 主表对齐
--
-- 背景：
--   customers 主表 uid 是 utf8mb4_unicode_ci
--   分表 batch_orders_N.user_uuid 是 utf8mb4_general_ci
--   JOIN 时报错：1267 - Illegal mix of collations
--
-- 本脚本遍历当前库中所有 batch_orders_N / batch_order_goods_N / batch_order_pay_N
-- 分表，把字符串列的排序规则改成 utf8mb4_unicode_ci（保留原长度和 nullable 属性）。
--
-- 注意：
--   1. MODIFY COLUMN 会重建该列的索引，大表可能耗时长，建议低峰执行
--   2. 字段类型保持 varchar(64) 与原定义一致
-- ======================================================================

-- 放宽 group_concat_max_len，避免分表多时被截断
SET SESSION group_concat_max_len = 1048576;

-- ==================== batch_orders_N ====================
-- 改：user_uuid, batch_uuid, owner_user
DROP PROCEDURE IF EXISTS fix_orders_collation;
DELIMITER $$
CREATE PROCEDURE fix_orders_collation()
fix_orders: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl VARCHAR(64);
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME REGEXP '^batch_orders_[0-9]+$'
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    SET @sql := '';
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO tbl;
        IF done = 1 THEN LEAVE read_loop; END IF;

        SET @sql := CONCAT(@sql,
            'ALTER TABLE `', tbl, '` ',
            'MODIFY COLUMN `user_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `batch_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `owner_user` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL; '
        );
    END LOOP;
    CLOSE cur;

    IF @sql = '' THEN
        SELECT '未找到 batch_orders_N 分表' AS msg;
        LEAVE fix_orders;
    END IF;

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    SELECT CONCAT('batch_orders_N 修改完成') AS msg;
END$$
DELIMITER ;

-- ==================== batch_order_goods_N ====================
-- 改：batch_order_uuid, batch_uuid, goods_uuid, owner_user, user_uuid
DROP PROCEDURE IF EXISTS fix_goods_collation;
DELIMITER $$
CREATE PROCEDURE fix_goods_collation()
fix_goods: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl VARCHAR(64);
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME REGEXP '^batch_order_goods_[0-9]+$'
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    SET @sql := '';
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO tbl;
        IF done = 1 THEN LEAVE read_loop; END IF;

        SET @sql := CONCAT(@sql,
            'ALTER TABLE `', tbl, '` ',
            'MODIFY COLUMN `batch_order_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `batch_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `goods_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `owner_user` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `user_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL; '
        );
    END LOOP;
    CLOSE cur;

    IF @sql = '' THEN
        SELECT '未找到 batch_order_goods_N 分表' AS msg;
        LEAVE fix_goods;
    END IF;

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    SELECT CONCAT('batch_order_goods_N 修改完成') AS msg;
END$$
DELIMITER ;

-- ==================== batch_order_pay_N ====================
-- 改：uid, customer_uuid, owner_user, batch_order_uuid
DROP PROCEDURE IF EXISTS fix_pay_collation;
DELIMITER $$
CREATE PROCEDURE fix_pay_collation()
fix_pay: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl VARCHAR(64);
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME REGEXP '^batch_order_pay_[0-9]+$'
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    SET @sql := '';
    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO tbl;
        IF done = 1 THEN LEAVE read_loop; END IF;

        SET @sql := CONCAT(@sql,
            'ALTER TABLE `', tbl, '` ',
            'MODIFY COLUMN `uid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL, ',
            'MODIFY COLUMN `customer_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `owner_user` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL, ',
            'MODIFY COLUMN `batch_order_uuid` VARCHAR(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL; '
        );
    END LOOP;
    CLOSE cur;

    IF @sql = '' THEN
        SELECT '未找到 batch_order_pay_N 分表' AS msg;
        LEAVE fix_pay;
    END IF;

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
    SELECT CONCAT('batch_order_pay_N 修改完成') AS msg;
END$$
DELIMITER ;

-- ==================== 执行 ====================
CALL fix_orders_collation();
CALL fix_goods_collation();
CALL fix_pay_collation();

-- ==================== 清理 ====================
DROP PROCEDURE IF EXISTS fix_orders_collation;
DROP PROCEDURE IF EXISTS fix_goods_collation;
DROP PROCEDURE IF EXISTS fix_pay_collation;

-- ==================== 验证 ====================
-- 检查分表的字符串列排序规则，应该全部为 utf8mb4_unicode_ci
SELECT
    TABLE_NAME,
    COLUMN_NAME,
    CHARACTER_SET_NAME,
    COLLATION_NAME
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME REGEXP '^batch_(orders|order_goods|order_pay)_[0-9]+$'
  AND CHARACTER_SET_NAME IS NOT NULL
ORDER BY TABLE_NAME, COLUMN_NAME;
