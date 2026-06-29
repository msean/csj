package logic

import (
	"app/global"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PaymentLogic struct {
	context   *gin.Context
	runtime   *global.RunTime
	OwnerUser string
}

func NewPaymentLogic(context *gin.Context) *PaymentLogic {
	return &PaymentLogic{
		context:   context,
		runtime:   global.Global,
		OwnerUser: common.GetUserUUID(context),
	}
}

// QuickPayList 快速还款列表 - 显示有赊欠的客户
func (logic *PaymentLogic) QuickPayList(req request.QuickPayListReq) (items []response.QuickPayListItem, total int64, err error) {
	// 使用 DAO 查询赊欠汇总
	var summaries []dao.CreditSummary
	if req.PageCount > 0 {
		offset := (req.Page - 1) * req.PageCount
		summaries, total, err = dao.PaymentDao.FindCreditSummariesWithPage(
			logic.runtime.DB, logic.OwnerUser, common.ValidOrder, offset, req.PageCount)
	} else {
		summaries, err = dao.PaymentDao.FindCreditSummaries(logic.runtime.DB, logic.OwnerUser, common.ValidOrder)
	}

	if err != nil {
		return
	}

	// 获取客户名称
	var userUUIDs []string
	for _, s := range summaries {
		userUUIDs = append(userUUIDs, s.UserUUID)
	}

	customerMap, _ := cache.CustomerCache.BatchCustomerFeildSet(userUUIDs, logic.OwnerUser)

	// 构建返回结果
	for _, s := range summaries {
		items = append(items, response.QuickPayListItem{
			CustomerUUID: s.UserUUID,
			CustomerName: customerMap[s.UserUUID].CustomerName,
			TotalCredit:  s.TotalCredit,
		})
	}

	return
}

// QuickPay 快捷还款 - 按时间顺序还款
func (logic *PaymentLogic) QuickPay(req request.QuickPayReq) (rsp response.QuickPayRsp, err error) {
	// 1. 使用 DAO 查询该客户所有未还清的订单，按时间排序
	orders, err := dao.PaymentDao.FindOrdersWithCredit(
		logic.runtime.DB, logic.OwnerUser, req.CustomerUUID, common.ValidOrder)
	if err != nil {
		return
	}

	if len(orders) == 0 {
		err = fmt.Errorf("该客户没有赊欠订单")
		return
	}

	// 2. 开始事务
	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. 按时间顺序还款
	remainingAmount := req.Amount
	var payDetails []model.PayDetail

	for i := range orders {
		if remainingAmount <= 0 {
			break
		}

		order := &orders[i]
		originalCredit := order.CreditAmount

		// 计算本次还款金额
		payAmount := remainingAmount
		if payAmount > order.CreditAmount {
			payAmount = order.CreditAmount
		}

		// 更新订单赊欠金额（可以为负）
		order.CreditAmount -= payAmount

		// 如果赊欠金额 <= 0，标记为已结清
		if order.CreditAmount <= 0 {
			order.Status = common.BatchOrderFinish
		}

		if err = tx.Save(order).Error; err != nil {
			return
		}

		remainingAmount -= payAmount

		// 记录还款详情
		payDetails = append(payDetails, model.PayDetail{
			OrderUUID: order.UID,
			Amount:    payAmount,
		})

		logic.runtime.Logger.Info("快捷还款",
			zap.String("orderUUID", order.UID),
			zap.Float64("originalCredit", originalCredit),
			zap.Float64("payAmount", payAmount),
			zap.Float64("newCredit", order.CreditAmount),
			zap.Int("status", order.Status))
	}

	// 4. 创建还款记录（使用分表 DAO）
	payDetailsJSON, _ := json.Marshal(payDetails)

	payRecord := model.BatchOrderPay{
		CustomerUUID:   req.CustomerUUID,
		OwnerUser:      logic.OwnerUser,
		BatchOrderUUID: "", // 快捷还款为空
		Amount:         req.Amount,
		PayType:        req.PayType,
		Remark:         req.Remark,
		PayDetails:     string(payDetailsJSON),
	}

	if err = dao.PaymentDao.CreatePay(tx, &payRecord); err != nil {
		return
	}

	// 5. 创建消息记录
	if err = logic.createMessage(tx, req.CustomerUUID, 1, "快捷还款",
		fmt.Sprintf("还款 %.2f 元", req.Amount), payRecord.UID); err != nil {
		return
	}

	// 6. 提交事务
	if err = tx.Commit().Error; err != nil {
		return
	}

	rsp.PayUUID = payRecord.UID
	return
}

// OrderPay 针对订单还款
func (logic *PaymentLogic) OrderPay(req request.OrderPayReq) (rsp response.OrderPayRsp, err error) {
	// 1. 使用 DAO 查询订单
	order, err := dao.PaymentDao.FindOrderByUID(logic.runtime.DB, logic.OwnerUser, req.OrderUUID)
	if err != nil {
		return
	}

	// 2. 开始事务
	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. 更新订单赊欠金额（可以为负）
	originalCredit := order.CreditAmount
	order.CreditAmount -= req.Amount

	// 如果赊欠金额 <= 0，标记为已结清
	if order.CreditAmount <= 0 {
		order.Status = common.BatchOrderFinish
	}

	if err = tx.Save(order).Error; err != nil {
		return
	}

	// 4. 创建还款记录（使用分表 DAO）
	payRecord := model.BatchOrderPay{
		CustomerUUID:   order.UserUUID,
		OwnerUser:      logic.OwnerUser,
		BatchOrderUUID: req.OrderUUID,
		Amount:         req.Amount,
		PayType:        req.PayType,
		Remark:         req.Remark,
		PayDetails:     "", // 针对订单还款为空
	}

	if err = dao.PaymentDao.CreatePay(tx, &payRecord); err != nil {
		return
	}

	// 5. 创建消息记录
	if err = logic.createMessage(tx, order.UserUUID, 1, "订单还款",
		fmt.Sprintf("订单 %s 还款 %.2f 元", req.OrderUUID[:8]+"...", req.Amount), payRecord.UID); err != nil {
		return
	}

	// 6. 提交事务
	if err = tx.Commit().Error; err != nil {
		return
	}

	logic.runtime.Logger.Info("订单还款",
		zap.String("orderUUID", req.OrderUUID),
		zap.Float64("originalCredit", originalCredit),
		zap.Float64("payAmount", req.Amount),
		zap.Float64("newCredit", order.CreditAmount),
		zap.Int("status", order.Status))

	rsp.PayUUID = payRecord.UID
	rsp.OrderUUID = req.OrderUUID
	rsp.OrderCredit = order.CreditAmount
	return
}

// PayHistory 还款历史
func (logic *PaymentLogic) PayHistory(req request.PayHistoryReq) (items []response.PayHistoryItem, total int64, err error) {
	// 使用 DAO 查询还款记录
	var offset int
	if req.PageCount > 0 {
		offset = (req.Page - 1) * req.PageCount
	}

	pays, total, err := dao.PaymentDao.FindPayRecords(
		logic.runtime.DB, logic.OwnerUser, req.CustomerUUID, offset, req.PageCount)
	if err != nil {
		return
	}

	// 获取客户名称
	customer, _ := dao.CustomerDao.FromUUID(logic.runtime.DB, req.CustomerUUID)

	for _, pay := range pays {
		items = append(items, response.PayHistoryItem{
			PayUUID:       pay.UID,
			CustomerUUID:  pay.CustomerUUID,
			CustomerName:  customer.Name,
			Amount:        pay.Amount,
			PayType:       pay.PayType,
			Remark:        pay.Remark,
			IsRevoked:     pay.IsRevoked,
			RevokedAt:     pay.RevokedAt,
			RevokedReason: pay.RevokedReason,
			CreatedAt:     pay.CreatedAt,
		})
	}

	return
}

// PayDetail 还款详情
func (logic *PaymentLogic) PayDetail(req request.PayDetailReq) (rsp response.PayDetailRsp, err error) {
	// 使用 DAO 查询还款记录
	pay, err := dao.PaymentDao.FindPayByUID(logic.runtime.DB, logic.OwnerUser, req.PayUUID)
	if err != nil {
		return
	}

	// 获取客户名称
	customer, _ := dao.CustomerDao.FromUUID(logic.runtime.DB, pay.CustomerUUID)

	rsp = response.PayDetailRsp{
		PayUUID:       pay.UID,
		CustomerUUID:  pay.CustomerUUID,
		CustomerName:  customer.Name,
		OrderUUID:     pay.BatchOrderUUID,
		Amount:        pay.Amount,
		PayType:       pay.PayType,
		Remark:        pay.Remark,
		PayDetails:    pay.PayDetails,
		IsRevoked:     pay.IsRevoked,
		RevokedAt:     pay.RevokedAt,
		RevokedReason: pay.RevokedReason,
		CreatedAt:     pay.CreatedAt,
	}

	return
}

// RevokePay 撤销还款
func (logic *PaymentLogic) RevokePay(req request.RevokePayReq) (err error) {
	// 1. 使用 DAO 查询还款记录
	pay, err := dao.PaymentDao.FindPayByUID(logic.runtime.DB, logic.OwnerUser, req.PayUUID)
	if err != nil {
		return
	}

	if pay.IsRevoked == 1 {
		err = fmt.Errorf("该还款记录已撤销")
		return
	}

	// 2. 开始事务
	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. 撤销还款 - 恢复订单赊欠金额
	if pay.BatchOrderUUID != "" {
		// 针对订单还款：直接恢复该订单
		order, err := dao.PaymentDao.FindOrderByUID(tx, logic.OwnerUser, pay.BatchOrderUUID)
		if err != nil {
			return err
		}

		order.CreditAmount += pay.Amount // 恢复原金额
		if err = tx.Save(order).Error; err != nil {
			return err
		}
	} else if pay.PayDetails != "" {
		// 快捷还款：按原还款详情恢复
		var payDetails []model.PayDetail
		if err = json.Unmarshal([]byte(pay.PayDetails), &payDetails); err != nil {
			return
		}

		for _, detail := range payDetails {
			order, err := dao.PaymentDao.FindOrderByUID(tx, logic.OwnerUser, detail.OrderUUID)
			if err != nil {
				logic.runtime.Logger.Warn("撤销还款 - 订单不存在", zap.String("orderUUID", detail.OrderUUID))
				continue
			}

			order.CreditAmount += detail.Amount
			if err = tx.Save(order).Error; err != nil {
				return err
			}

			logic.runtime.Logger.Info("撤销还款 - 恢复订单",
				zap.String("orderUUID", detail.OrderUUID),
				zap.Float64("amount", detail.Amount))
		}
	}

	// 4. 更新还款记录状态
	now := time.Now()
	pay.IsRevoked = 1
	pay.RevokedAt = &now
	pay.RevokedReason = req.Reason

	if err = tx.Save(pay).Error; err != nil {
		return
	}

	// 5. 创建消息记录
	if err = logic.createMessage(tx, pay.CustomerUUID, 2, "撤销还款",
		fmt.Sprintf("撤销还款 %.2f 元，原因：%s", pay.Amount, req.Reason), pay.UID); err != nil {
		return
	}

	// 6. 提交事务
	if err = tx.Commit().Error; err != nil {
		return
	}

	logic.runtime.Logger.Info("撤销还款",
		zap.String("payUUID", pay.UID),
		zap.Float64("amount", pay.Amount),
		zap.String("reason", req.Reason))

	return
}

// MessageList 消息列表
func (logic *PaymentLogic) MessageList(req request.MessageListReq) (items []response.MessageItem, summary response.MessageSummary, err error) {
	// 使用 DAO 统计消息
	totalCount, unreadCount, err := dao.PaymentDao.CountMessages(logic.runtime.DB, logic.OwnerUser)
	if err != nil {
		return
	}

	summary.TotalCount = int(totalCount)
	summary.UnreadCount = int(unreadCount)

	// 分页查询消息
	var offset int
	if req.PageCount > 0 {
		offset = (req.Page - 1) * req.PageCount
	}

	messages, err := dao.PaymentDao.FindMessages(logic.runtime.DB, logic.OwnerUser, offset, req.PageCount)
	if err != nil {
		return
	}

	// 获取客户名称
	customerUUIDs := make(map[string]bool)
	for _, msg := range messages {
		customerUUIDs[msg.CustomerUUID] = true
	}

	var uuids []string
	for uuid := range customerUUIDs {
		uuids = append(uuids, uuid)
	}

	customerMap, _ := cache.CustomerCache.BatchCustomerFeildSet(uuids, logic.OwnerUser)

	for _, msg := range messages {
		items = append(items, response.MessageItem{
			MessageUUID:  msg.UID,
			Type:         msg.Type,
			Event:        msg.Event,
			Content:      msg.Content,
			CustomerUUID: msg.CustomerUUID,
			CustomerName: customerMap[msg.CustomerUUID].CustomerName,
			IsRead:       msg.IsRead,
			RelatedUUID:  msg.RelatedUUID,
			CreatedAt:    msg.CreatedAt,
		})
	}

	return
}

// createMessage 创建消息记录
func (logic *PaymentLogic) createMessage(tx *gorm.DB, customerUUID string, msgType int, event, content, relatedUUID string) error {
	message := model.MessageCenter{
		OwnerUser:    logic.OwnerUser,
		CustomerUUID: customerUUID,
		Type:         msgType,
		Event:        event,
		Content:      content,
		RelatedUUID:  relatedUUID,
		IsRead:       0,
	}
	return dao.PaymentDao.CreateMessage(tx, &message)
}
