-- ======================================================================
-- 将 customer.debt 迁移到 customer.total_credit（Credit字段），并清空 debt / paid_debt
--
-- 背景：
--   Credit (= total_credit) 表示客户当前总欠款，与 Debt 是一对镜像：
--     - 建客户时 Credit = Debt
--     - 改 Debt 时，CustomerDao.Update 按 delta 等量同步 Credit
--     - 还款时 DecrDebt(actualPaid) 同时减少 Credit
--     - 撤销还款时 IncrDebt(restoreDebt) 同步增加 Credit
--
--   此脚本用于"重置 Debt 归零"：历史上累计到 Debt 的值已经在 Credit 中体现，
--   执行后 Debt / PaidDebt 都归零，后续 Debt 将作为"新录入的客户原始欠款"重新累积。
--
-- 注意：
--   1. 本操作会改变后续快捷还款/订单还款的分配顺序（PaidDebt=Debt=0 时跳过原始债务，直接抵订单）
--   2. 执行前建议对 customers 表做一次备份
--   3. 如需保留 paid_debt 历史值，把最后一条 paid_debt=0 的 UPDATE 注释掉
-- ======================================================================

-- 可选：开始事务（生产环境建议打开，便于出错时 ROLLBACK）
-- START TRANSACTION;

-- Step 1. 将 debt 的值累加到 total_credit（Credit 字段）
--         用 debt 加而不是直接覆盖，避免 Credit 中已经通过还款/订单累积的增量丢失
UPDATE `customers`
SET `total_credit` = `total_credit` + `debt`
WHERE `debt` <> 0;

-- Step 2. 将 debt 清零
UPDATE `customers`
SET `debt` = 0
WHERE `debt` <> 0;

-- Step 3. 将 paid_debt 清零（与 debt 配对；如需保留历史则注释掉）
UPDATE `customers`
SET `paid_debt` = 0
WHERE `paid_debt` <> 0;

-- 可选：提交事务
-- COMMIT;

-- 验证：统计归零后 debt / paid_debt 非零的客户数（应为 0）
SELECT
    COUNT(*)                                             AS total,
    SUM(CASE WHEN `debt` <> 0 THEN 1 ELSE 0 END)        AS debt_non_zero,
    SUM(CASE WHEN `paid_debt` <> 0 THEN 1 ELSE 0 END)   AS paid_debt_non_zero,
    SUM(`total_credit`)                                  AS total_credit_sum
FROM `customers`;
