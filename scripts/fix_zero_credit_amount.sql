-- ======================================================================
-- 修复 BatchOrder.CreditAmount 为 0 的记录：用 TotalAmount 赋值
-- 公式：IF credit_amount = 0 THEN credit_amount = amount
--
-- 说明：用存储过程自动遍历所有 batch_orders_N 分表
-- ======================================================================

DROP PROCEDURE IF EXISTS fix_zero_credit_amount;
DELIMITER $$

CREATE PROCEDURE fix_zero_credit_amount()
fix_zero_credit: BEGIN
    DECLARE done INT DEFAULT 0;
    DECLARE tbl_name VARCHAR(64);
    DECLARE affected_rows INT DEFAULT 0;
    DECLARE total_rows INT DEFAULT 0;

    DECLARE cur CURSOR FOR
        SELECT TABLE_NAME
        FROM information_schema.TABLES
        WHERE TABLE_SCHEMA = DATABASE()
          AND TABLE_NAME LIKE 'batch\_orders\_%'
          AND TABLE_NAME REGEXP '^batch_orders_[0-9]+$'
        ORDER BY TABLE_NAME;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = 1;

    OPEN cur;
    tbl_loop: LOOP
        FETCH cur INTO tbl_name;
        IF done = 1 THEN
            LEAVE tbl_loop;
        END IF;

        -- 对每个分表执行更新：credit_amount 为 0 时用 amount 赋值
        SET @sql = CONCAT(
            'UPDATE `', tbl_name, '` ',
            'SET `credit_amount` = `amount` ',
            'WHERE `credit_amount` = 0 AND `amount` > 0 AND `deleted_at` IS NULL'
        );
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        SET affected_rows = ROW_COUNT();
        SET total_rows = total_rows + affected_rows;
        DEALLOCATE PREPARE stmt;

        -- 可选：打印每个表的影响行数
        SELECT CONCAT('表 ', tbl_name, ' 更新了 ', affected_rows, ' 行') AS log;
    END LOOP;
    CLOSE cur;

    SELECT CONCAT('总计更新了 ', total_rows, ' 行') AS result;
END$$

DELIMITER ;

-- 执行存储过程
CALL fix_zero_credit_amount();

-- 清理
DROP PROCEDURE IF EXISTS fix_zero_credit_amount;
