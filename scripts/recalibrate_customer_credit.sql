-- ======================================================================
-- 校准 Customer.Credit 和 Customer.TotalAmount 字段
-- 公式：
--   Credit = Debt + SUM(订单 amount WHERE status=1) - SUM(还款 amount+discount WHERE is_revoked=0)
--   TotalAmount = SUM(所有订单 amount)
--
--   有效订单：status = 1（仅赊欠单），deleted_at IS NULL
--   有效还款：is_revoked = 0，deleted_at IS NULL
--   TotalAmount 统计：所有 deleted_at IS NULL 的订单（不限 status）
--
-- 说明：用存储过程自动遍历分表拼接 UNION ALL
-- ======================================================================

SET SESSION group_concat_max_len = 1048576;

DROP PROCEDURE IF EXISTS recalibrate_customer_credit;
DELIMITER $$

CREATE PROCEDURE recalibrate_customer_credit()
recalibrate_customer_credit: BEGIN
    -- 所有 DECLARE 必须放在最前面
    DECLARE done INT DEFAULT 0;
    DECLARE tbl_name VARCHAR(64);
    DECLARE order_union LONGTEXT DEFAULT '';
    DECLARE pay_union LONGTEXT DEFAULT '';
    DECLARE full_sql LONGTEXT;

    DECLARE cur_order CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME LIKE 'batch\_orders\_%'
          AND TABLE_NAME REGEXP '^batch_orders_[0-9]+$'
        ORDER BY TABLE_NAME;

    DECLARE cur_pay CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME LIKE 'batch\_order\_pay\_%'
          AND TABLE_NAME REGEXP '^batch_order_pay_[0-9]+$'
        ORDER BY TABLE_NAME;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    -- ===== 1. 拼接订单分表 UNION ALL =====
    -- 同时输出 credit_order_amount（status=1 的金额之和）和 total_order_amount（所有订单金额之和）
    OPEN cur_order;
    order_loop: LOOP
        FETCH cur_order INTO tbl_name;
        IF done = 1 THEN
            LEAVE order_loop;
        END IF;

        IF order_union = '' THEN
            SET order_union = CONCAT(
                'SELECT user_uuid, ',
                'SUM(CASE WHEN `status` = 1 THEN `amount` ELSE 0 END) AS credit_order_amount, ',
                'SUM(`amount`) AS total_order_amount ',
                'FROM `', tbl_name, '` ',
                'WHERE `deleted_at` IS NULL AND `user_uuid` <> '''' ',
                'GROUP BY user_uuid'
            );
        ELSE
            SET order_union = CONCAT(
                order_union,
                ' UNION ALL ',
                'SELECT user_uuid, ',
                'SUM(CASE WHEN `status` = 1 THEN `amount` ELSE 0 END) AS credit_order_amount, ',
                'SUM(`amount`) AS total_order_amount ',
                'FROM `', tbl_name, '` ',
                'WHERE `deleted_at` IS NULL AND `user_uuid` <> '''' ',
                'GROUP BY user_uuid'
            );
        END IF;
    END LOOP;
    CLOSE cur_order;

    -- ===== 2. 拼接还款分表 UNION ALL =====
    SET done = 0;
    OPEN cur_pay;
    pay_loop: LOOP
        FETCH cur_pay INTO tbl_name;
        IF done = 1 THEN
            LEAVE pay_loop;
        END IF;

        IF pay_union = '' THEN
            SET pay_union = CONCAT(
                'SELECT customer_uuid, (`amount` + `discount`) AS pay_amount ',
                'FROM `', tbl_name, '` ',
                'WHERE `is_revoked` = 0 AND `deleted_at` IS NULL AND `customer_uuid` <> '''''
            );
        ELSE
            SET pay_union = CONCAT(
                pay_union,
                ' UNION ALL ',
                'SELECT customer_uuid, (`amount` + `discount`) AS pay_amount ',
                'FROM `', tbl_name, '` ',
                'WHERE `is_revoked` = 0 AND `deleted_at` IS NULL AND `customer_uuid` <> '''''
            );
        END IF;
    END LOOP;
    CLOSE cur_pay;

    -- 没有订单分表也没有还款分表，退出
    IF order_union = '' AND pay_union = '' THEN
        SELECT '未找到分表，退出' AS msg;
        LEAVE recalibrate_customer_credit;
    END IF;

    -- ===== 3. 构造校准 SQL =====
    -- 同时更新 credit 和 total_amount
    SET full_sql = CONCAT(
        'UPDATE `customers` c ',
        'LEFT JOIN (',
            'SELECT user_uuid, ',
            'COALESCE(SUM(credit_order_amount), 0) AS order_credit_total, ',
            'COALESCE(SUM(total_order_amount), 0) AS order_total ',
            'FROM (', IFNULL(order_union, 'SELECT '''' AS user_uuid, 0 AS credit_order_amount, 0 AS total_order_amount WHERE 1=0'), ') t ',
            'GROUP BY user_uuid',
        ') o ON o.user_uuid = c.uid COLLATE utf8mb4_unicode_ci ',
        'LEFT JOIN (',
            'SELECT customer_uuid, COALESCE(SUM(pay_amount), 0) AS pay_total ',
            'FROM (', IFNULL(pay_union, 'SELECT '''' AS customer_uuid, 0 AS pay_amount WHERE 1=0'), ') p ',
            'GROUP BY customer_uuid',
        ') py ON py.customer_uuid = c.uid COLLATE utf8mb4_unicode_ci ',
        'SET c.credit = c.debt + COALESCE(o.order_credit_total, 0) - COALESCE(py.pay_total, 0), ',
        '    c.total_amount = COALESCE(o.order_total, 0)'
    );

    SET @sql = full_sql;
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SELECT '校准完成' AS msg;
END$$

DELIMITER ;

-- 执行存储过程
CALL recalibrate_customer_credit();

-- 清理
DROP PROCEDURE IF EXISTS recalibrate_customer_credit;
