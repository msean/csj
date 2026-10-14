-- ======================================================================
-- 给 Customer.Debt 字段赋值（存量数据回填）
-- 规则：该客户在所有 batch_orders 分表中，status IN (1,2,3) 的 credit_amount 之和
--   1 = BatchOrderTemp      记账单
--   2 = BatchOrderedCredit   赊欠单
--   3 = BatchOrderFinish     已结算
-- 关联：customers.uid = batch_orders_N.user_uuid
--
-- 说明：用存储过程自动发现当前库中已存在的 batch_orders_N 分表并拼接 UNION，
--       避免硬编码表名，缺失的分表会自动跳过。
-- 注意：需放宽 group_concat_max_len，否则分表多时会被截断。
-- ======================================================================

-- 放宽 GROUP_CONCAT 长度，避免 50 个分表拼接被截断
SET SESSION group_concat_max_len = 1048576;

DROP PROCEDURE IF EXISTS init_customer_debt;
DELIMITER $$

CREATE PROCEDURE init_customer_debt()
init_customer_debt: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl_name VARCHAR(64);
    DECLARE union_sql LONGTEXT DEFAULT '';
    DECLARE full_sql LONGTEXT;
    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME LIKE 'batch\_orders\_%'
          AND TABLE_NAME REGEXP '^batch_orders_[0-9]+$'
        ORDER BY TABLE_NAME;
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    OPEN cur;
    read_loop: LOOP
        FETCH cur INTO tbl_name;
        IF done = 1 THEN
            LEAVE read_loop;
        END IF;

        IF union_sql = '' THEN
            SET union_sql = CONCAT(
                'SELECT user_uuid, credit_amount FROM `', tbl_name, '` ',
                'WHERE `status` IN (1,2,3) AND `deleted_at` IS NULL AND `user_uuid` <> '''''
            );
        ELSE
            SET union_sql = CONCAT(
                union_sql,
                ' UNION ALL ',
                'SELECT user_uuid, credit_amount FROM `', tbl_name, '` ',
                'WHERE `status` IN (1,2,3) AND `deleted_at` IS NULL AND `user_uuid` <> '''''
            );
        END IF;
    END LOOP;
    CLOSE cur;

    -- 没有分表，直接退出
    IF union_sql = '' THEN
        SELECT '未找到 batch_orders_N 分表，退出' AS msg;
        LEAVE init_customer_debt;
    END IF;

    -- 先把所有客户 credit 清零（避免历史脏数据残留），再用聚合结果覆盖
    -- 注意：只清零有对应订单的客户由下面的 JOIN 限定，这里全量清零是为了
    --       让"曾经有订单但分表已无记录"的客户也能归零
    UPDATE `customers` SET credit = 0 WHERE credit != 0;

    SET full_sql = CONCAT(
        'UPDATE `customers` c ',
        'JOIN (',
            'SELECT user_uuid, SUM(credit_amount) AS total_debt ',
            'FROM (', union_sql, ') t ',
            'GROUP BY user_uuid',
        ') agg ON agg.user_uuid = c.uid COLLATE utf8mb4_unicode_ci ',
        'SET c.credit = agg.total_debt'
    );

    SET @sql = full_sql;
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SELECT CONCAT('初始化完成，共处理分表') AS msg;
END$$

DELIMITER ;

-- 执行存储过程
CALL init_customer_debt();

-- 清理
DROP PROCEDURE IF EXISTS init_customer_debt;
