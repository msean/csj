-- ======================================================================
-- 删除分表架构中主表上残留的外键约束
--
-- 背景：
--   项目已将 batch_orders / batch_order_goods / batch_order_pay 改为分表
--   （batch_orders_N / batch_order_goods_N / batch_order_pay_N）。
--   主表上由 GORM AutoMigrate 创建的外键约束在分表场景下不再需要，
--   且会阻止向主表插入/更新数据（Error 1452 外键约束失败）。
--
-- 说明：
--   本脚本会自动查询并删除 batch_order_goods / batch_order_pay / batch_orders
--   三张主表上的所有外键约束，可重复执行（幂等）。
--
-- 用法：
--   mysql -u root -p caishuji < scripts/drop_shard_foreign_keys.sql
--   或在 MySQL 客户端直接粘贴执行。
-- ======================================================================

-- 放宽 GROUP_CONCAT 长度
SET SESSION group_concat_max_len = 1048576;

-- 存储过程：动态删除指定表上的所有外键约束
DROP PROCEDURE IF EXISTS drop_shard_foreign_keys;
DELIMITER $$

CREATE PROCEDURE drop_shard_foreign_keys()
proc_main: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE fk_name VARCHAR(64);
    DECLARE tbl_name VARCHAR(64);
    DECLARE drop_sql TEXT;
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME, CONSTRAINT_NAME
        FROM information_schema.TABLE_CONSTRAINTS
        WHERE CONSTRAINT_SCHEMA = DATABASE()
          AND CONSTRAINT_TYPE = 'FOREIGN KEY'
          AND TABLE_NAME IN ('batch_order_goods', 'batch_order_pay', 'batch_orders')
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO tbl_name, fk_name;
        IF done = 1 THEN
            LEAVE read_loop;
        END IF;

        SET @drop_sql = CONCAT('ALTER TABLE `', tbl_name, '` DROP FOREIGN KEY `', fk_name, '`');
        PREPARE stmt FROM @drop_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;

        SELECT CONCAT('已删除外键: ', tbl_name, '.', fk_name) AS msg;
    END LOOP;
    CLOSE cur;

    SELECT '外键约束清理完成' AS msg;
END$$

DELIMITER ;

-- 执行
CALL drop_shard_foreign_keys();

-- 清理
DROP PROCEDURE IF EXISTS drop_shard_foreign_keys;
