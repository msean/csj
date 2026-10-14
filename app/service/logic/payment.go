package logic

import (
	"app/global"
	"app/pkg/utils"
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
func (logic *PaymentLogic) QuickPayList(req request.QuickPayListReq) (items []response.QuickPayListItem, total int64, totalDebt float64, err error) {
	// 直接从 Customer 表查 debt != 0 的客户，无需走订单分表聚合
	offset := 0
	if req.PageCount > 0 {
		offset = (req.Page - 1) * req.PageCount
	}
	list, total, totalDebt, err := dao.CustomerDao.FindDebtCustomers(logic.runtime.DB, logic.OwnerUser, req.Search, offset, req.PageCount)
	if err != nil {
		return
	}

	items = make([]response.QuickPayListItem, 0, len(list))
	for _, c := range list {
		items = append(items, response.QuickPayListItem{
			CustomerUUID: c.UID,
			CustomerName: c.Name,
			TotalCredit:  c.Debt,
		})
	}
	return
}

// QuickPay 快捷还款 - 按时间顺序还款
func (logic *PaymentLogic) QuickPay(req request.QuickPayReq) (rsp response.QuickPayRsp, err error) {
	c, err := dao.CustomerDao.FromUUID(logic.runtime.DB, req.CustomerUUID)
	if err != nil {
		return
	}
	// 记录还款前客户赊欠总额（订单级 CreditAmount 之和）
	beforeCredit := c.Credit
	// 2. 开始事务

	orders, err := dao.OrderDao.FindOrdersWithCredit(
		logic.runtime.DB, logic.OwnerUser, req.CustomerUUID, common.UnFinalStatusOrder)
	if err != nil {
		return
	}

	tx := logic.runtime.DB.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 3. 先判断客户是否存在未清的"原始债务"（Debt > PaidDebt）
	//    还款优先级：先还原始债务（PaidDebt）→ 还完后再抵扣订单赊欠

	pendingDebt := c.Debt - c.PaidDebt
	if pendingDebt < 0 {
		pendingDebt = 0
	}

	// 还款总可用资金 = 用户还款 + 优惠
	remainingAmount := req.Amount + req.DisCount
	var payDetails []model.PayDetail

	// 3a. 优先偿还原始债务
	debtPaidAmount := 0.0 // 累计用于偿还原始债务的金额
	if pendingDebt > 0 && remainingAmount > 0 {
		debtPaidAmount = remainingAmount
		if debtPaidAmount > pendingDebt {
			debtPaidAmount = pendingDebt
		}
		remainingAmount -= debtPaidAmount

		if err = dao.CustomerDao.Pay(tx, req.CustomerUUID, debtPaidAmount); err != nil {
			return
		}

		payDetails = append(payDetails, model.PayDetail{
			Type:      common.PayDetailTypeDebt,
			OrderUUID: "",
			Amount:    debtPaidAmount,
		})

		logic.runtime.Logger.Info("快捷还款-先还原始债务",
			zap.String("customerUUID", req.CustomerUUID),
			zap.Float64("pendingDebt", pendingDebt),
			zap.Float64("debtPaidAmount", debtPaidAmount),
			zap.Float64("remainingAfterDebt", remainingAmount))
	}

	// 3b. 若还有剩余金额，再按时间顺序抵扣订单赊欠
	var hasUnpaid bool
	// 1. 使用 DAO 查询该客户所有未还清的订单，按时间排序

	if remainingAmount > 0 {
		for i := range orders {
			if remainingAmount <= 0 {
				break
			}

			order := &orders[i]
			originalCredit := order.CreditAmount

			payAmount := remainingAmount
			if payAmount > order.CreditAmount {
				payAmount = order.CreditAmount
			}

			order.CreditAmount -= payAmount

			if order.CreditAmount <= 0 {
				order.Status = common.BatchOrderFinish
			}

			if err = dao.OrderDao.UpdateOrderCreditStatus(tx, logic.OwnerUser, order.UID, order.CreditAmount, order.Status); err != nil {
				return
			}

			remainingAmount -= payAmount

			payDetails = append(payDetails, model.PayDetail{
				Type:      common.PayDetailTypeOrder,
				OrderUUID: order.UID,
				Amount:    payAmount,
			})

			logic.runtime.Logger.Info("快捷还款-抵扣订单",
				zap.String("orderUUID", order.UID),
				zap.Float64("originalCredit", originalCredit),
				zap.Float64("payAmount", payAmount),
				zap.Float64("newCredit", order.CreditAmount),
				zap.Int("status", order.Status))
		}

		// 判断是否仍有未还清的订单赊欠
		var afterCredit float64
		for i := range orders {
			if orders[i].CreditAmount > 0 {
				afterCredit += orders[i].CreditAmount
			}
		}
		hasUnpaid = afterCredit > 0
	} else {
		// 还款金额全部用于还原始债务，订单赊欠一笔都没抵 —— 肯定未还清
		hasUnpaid = true
	}

	// 4. 创建还款记录（使用分表 DAO）
	payDetailsJSON, _ := json.Marshal(payDetails)

	payRecord := model.BatchOrderPay{
		CustomerUUID:   req.CustomerUUID,
		OwnerUser:      logic.OwnerUser,
		BatchOrderUUID: "", // 快捷还款为空
		Amount:         req.Amount,
		Discount:       req.DisCount,
		PayType:        req.PayType,
		PayDetails:     string(payDetailsJSON),
		BeforeCredit:   beforeCredit,
	}

	if err = dao.PaymentDao.CreatePay(tx, &payRecord); err != nil {
		return
	}

	// 5. 扣减客户总欠款 Customer.Debt
	//    注意：实际抵扣金额 = 还款金额+优惠 - 未用完的剩余金额
	//    包含用于原始债务和订单赊欠的总和（Debt 是总欠款，扣减一致）
	actualPaid := req.Amount + req.DisCount - remainingAmount
	if actualPaid > 0 {
		if err = dao.CustomerDao.DecrDebt(tx, req.CustomerUUID, actualPaid); err != nil {
			return
		}
	}

	// 6. 创建消息记录
	if err = logic.createMessage(tx, req.CustomerUUID, 1, "快捷还款",
		fmt.Sprintf("还款 %.2f 元", req.Amount), payRecord.UID); err != nil {
		return
	}

	// 7. 提交事务
	if err = tx.Commit().Error; err != nil {
		return
	}

	rsp.PayUUID = payRecord.UID
	rsp.Discount = req.DisCount
	rsp.Amount = req.Amount
	rsp.NeedPop = hasUnpaid
	return
}

// OrderPay 针对订单还款
func (logic *PaymentLogic) OrderPay(req request.OrderPayReq) (rsp response.OrderPayRsp, err error) {
	// 1. 使用 DAO 查询订单
	order, err := dao.OrderDao.FindByUID(logic.runtime.DB, logic.OwnerUser, req.OrderUUID)
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

	// 3. 先判断客户是否存在未清的"原始债务"（Debt > PaidDebt）
	//    还款优先级：先还原始债务 → 再还订单
	c, err := dao.CustomerDao.FromUUID(tx, order.UserUUID)
	if err != nil {
		return
	}
	pendingDebt := c.Debt - c.PaidDebt
	if pendingDebt < 0 {
		pendingDebt = 0
	}

	remainingAmount := req.Amount
	var payDetails []model.PayDetail

	// 3a. 优先偿还原始债务
	debtPaidAmount := 0.0
	if pendingDebt > 0 && remainingAmount > 0 {
		debtPaidAmount = remainingAmount
		if debtPaidAmount > pendingDebt {
			debtPaidAmount = pendingDebt
		}
		remainingAmount -= debtPaidAmount

		if err = dao.CustomerDao.Pay(tx, order.UserUUID, debtPaidAmount); err != nil {
			return
		}

		payDetails = append(payDetails, model.PayDetail{
			Type:      common.PayDetailTypeDebt,
			OrderUUID: "",
			Amount:    debtPaidAmount,
		})
	}

	// 3b. 剩余金额抵扣订单赊欠
	originalCredit := order.CreditAmount
	orderCreditPaid := 0.0
	if remainingAmount > 0 {
		orderCreditPaid = remainingAmount
		if orderCreditPaid > order.CreditAmount {
			orderCreditPaid = order.CreditAmount
		}
		remainingAmount -= orderCreditPaid
		order.CreditAmount -= orderCreditPaid

		if order.CreditAmount <= 0 {
			order.Status = common.BatchOrderFinish
		}

		if err = dao.OrderDao.UpdateOrderCreditStatus(tx, logic.OwnerUser, order.UID, order.CreditAmount, order.Status); err != nil {
			return
		}

		payDetails = append(payDetails, model.PayDetail{
			Type:      common.PayDetailTypeOrder,
			OrderUUID: order.UID,
			Amount:    orderCreditPaid,
		})
	}

	// 4. 创建还款记录（使用分表 DAO）
	payDetailsJSON, _ := json.Marshal(payDetails)
	payRecord := model.BatchOrderPay{
		CustomerUUID:   order.UserUUID,
		OwnerUser:      logic.OwnerUser,
		BatchOrderUUID: req.OrderUUID,
		Amount:         req.Amount,
		PayType:        req.PayType,
		PayDetails:     string(payDetailsJSON),
	}

	if err = dao.PaymentDao.CreatePay(tx, &payRecord); err != nil {
		return
	}

	// 5. 扣减客户总欠款 Customer.Debt
	//    实际扣减 = 总还款 - 未用完的剩余金额 = 原始债务抵扣 + 订单赊欠抵扣
	actualPaid := req.Amount - remainingAmount
	if actualPaid > 0 {
		if err = dao.CustomerDao.DecrDebt(tx, order.UserUUID, actualPaid); err != nil {
			return
		}
	}

	// 6. 创建消息记录
	if err = logic.createMessage(tx, order.UserUUID, 1, "订单还款",
		fmt.Sprintf("订单 %s 还款 %.2f 元", req.OrderUUID[:8]+"...", req.Amount), payRecord.UID); err != nil {
		return
	}

	// 7. 提交事务
	if err = tx.Commit().Error; err != nil {
		return
	}

	logic.runtime.Logger.Info("订单还款",
		zap.String("orderUUID", req.OrderUUID),
		zap.Float64("originalCredit", originalCredit),
		zap.Float64("debtPaidAmount", debtPaidAmount),
		zap.Float64("orderCreditPaid", orderCreditPaid),
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
		PayDetails:    pay.PayDetails,
		IsRevoked:     pay.IsRevoked,
		RevokedAt:     pay.RevokedAt,
		RevokedReason: pay.RevokedReason,
		Discount:      int(pay.Discount),
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

	// 3. 撤销还款 - 恢复原始债务/订单赊欠 与 客户总欠款
	// 按 PayDetail.Type 分别还原：type=1 还原 PaidDebt，type=2 还原订单 CreditAmount
	// Customer.Debt 的恢复总和 = 所有 PayDetail.Amount 之和 + Discount
	//    （等价于还款时实际抵扣的 Amount + Discount）
	restoreDebt := pay.Amount + pay.Discount

	if pay.BatchOrderUUID != "" {
		// 针对订单还款：从 PayDetails 中判断还款分配
		if pay.PayDetails != "" {
			var payDetails []model.PayDetail
			if err = json.Unmarshal([]byte(pay.PayDetails), &payDetails); err != nil {
				return
			}

			debtRestore := 0.0  // 需要还原 PaidDebt 的累计
			orderRestore := 0.0 // 需要还原订单 CreditAmount 的累计
			for _, detail := range payDetails {
				switch detail.Type {
				case common.PayDetailTypeDebt:
					// 还原 PaidDebt
					debtRestore += detail.Amount
					if err = dao.CustomerDao.DecrPaidDebt(tx, pay.CustomerUUID, detail.Amount); err != nil {
						return err
					}
				case common.PayDetailTypeOrder:
					// 还原订单赊欠
					orderRestore += detail.Amount
				}
			}

			// 还原订单 CreditAmount（订单还款必定有 BatchOrderUUID，直接用它查）
			if orderRestore > 0 {
				order, err := dao.OrderDao.FindByUID(tx, logic.OwnerUser, pay.BatchOrderUUID)
				if err != nil {
					return err
				}
				order.CreditAmount += orderRestore
				// 恢复赊欠后状态如果不是结清就还原为记账单（保持与还款逻辑一致）
				if order.CreditAmount > 0 {
					order.Status = common.BatchOrderTemp
				}
				if err = dao.OrderDao.UpdateOrderCreditStatus(tx, logic.OwnerUser, order.UID, order.CreditAmount, order.Status); err != nil {
					return err
				}
			}

			logic.runtime.Logger.Info("撤销还款-订单还款",
				zap.String("payUUID", pay.UID),
				zap.String("orderUUID", pay.BatchOrderUUID),
				zap.Float64("debtRestore", debtRestore),
				zap.Float64("orderRestore", orderRestore))
		} else {
			// 老数据：PayDetails 为空（旧版本记录），按原始逻辑整笔恢复订单和客户欠款
			order, err := dao.OrderDao.FindByUID(tx, logic.OwnerUser, pay.BatchOrderUUID)
			if err != nil {
				return err
			}
			order.CreditAmount += restoreDebt
			if err = dao.OrderDao.UpdateOrderCreditStatus(tx, logic.OwnerUser, order.UID, order.CreditAmount, order.Status); err != nil {
				return err
			}
		}
		// 统一恢复客户总欠款 Customer.Debt
		if err = dao.CustomerDao.IncrDebt(tx, pay.CustomerUUID, restoreDebt); err != nil {
			return err
		}
	} else if pay.PayDetails != "" {
		// 快捷还款：按 PayDetail.Type 分别还原 PaidDebt / 订单
		var payDetails []model.PayDetail
		if err = json.Unmarshal([]byte(pay.PayDetails), &payDetails); err != nil {
			return
		}

		for _, detail := range payDetails {
			switch detail.Type {
			case common.PayDetailTypeDebt:
				// 还原 PaidDebt
				if err = dao.CustomerDao.DecrPaidDebt(tx, pay.CustomerUUID, detail.Amount); err != nil {
					return err
				}
				logic.runtime.Logger.Info("撤销还款-还原原始债务",
					zap.String("customerUUID", pay.CustomerUUID),
					zap.Float64("amount", detail.Amount))

			case common.PayDetailTypeOrder:
				// 还原订单赊欠
				order, err := dao.OrderDao.FindByUID(tx, logic.OwnerUser, detail.OrderUUID)
				if err != nil {
					logic.runtime.Logger.Warn("撤销还款 - 订单不存在", zap.String("orderUUID", detail.OrderUUID))
					continue
				}
				order.CreditAmount += detail.Amount
				if order.CreditAmount > 0 {
					order.Status = common.BatchOrderTemp
				}
				if err = dao.OrderDao.UpdateOrderCreditStatus(tx, logic.OwnerUser, order.UID, order.CreditAmount, order.Status); err != nil {
					return err
				}
				logic.runtime.Logger.Info("撤销还款-还原订单",
					zap.String("orderUUID", detail.OrderUUID),
					zap.Float64("amount", detail.Amount))
			}
		}

		// 快捷还款只针对单个客户，统一恢复客户总欠款 Customer.Debt
		if err = dao.CustomerDao.IncrDebt(tx, pay.CustomerUUID, restoreDebt); err != nil {
			return err
		}
	}

	// 4. 更新还款记录状态
	now := time.Now()
	if err = dao.PaymentDao.UpdatePayRevokeStatus(tx, logic.OwnerUser, pay.UID, now, req.Reason); err != nil {
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

// MessageRead 标记消息为已读
func (logic *PaymentLogic) MessageRead(req request.MessageReadReq) error {
	return dao.PaymentDao.MarkMessageRead(logic.runtime.DB, logic.OwnerUser, req.UUID)
}

// GenerateShareHTML builds the payment receipt HTML and returns the HTML string.
func (logic *PaymentLogic) GenerateShareHTML(req request.PaymentShareReq) (html string, err error) {
	// 1. 查询还款记录
	var pay *model.BatchOrderPay
	if req.PayUUID != "" {
		pay, err = dao.PaymentDao.FindPayByUID(logic.runtime.DB, logic.OwnerUser, req.PayUUID)
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("PayUUID不能为空")
		return
	}

	// 2. 查询客户信息
	customer, err := dao.CustomerDao.FromUUID(logic.runtime.DB, pay.CustomerUUID)
	if err != nil {
		logic.runtime.Logger.Error("GenerateShareHTML customer not found", zap.Error(err))
		return
	}

	// 3. 计算展示数据
	preCredit := pay.BeforeCredit
	paidAmount := pay.Amount
	discount := pay.Discount
	currentAmount := preCredit - paidAmount - discount
	if currentAmount < 0 {
		currentAmount = 0
	}

	// 支付方式文案
	payTypeStr := ""
	switch pay.PayType {
	case common.PayTypeWx:
		payTypeStr = "微信"
	case common.PayTypeZFB:
		payTypeStr = "支付宝"
	case common.PayTypeBank:
		payTypeStr = "银行"
	case common.PayTypeCash:
		payTypeStr = "现金"
	default:
		payTypeStr = ""
	}
	paidAmountStr := utils.FormatFloat(paidAmount)
	if payTypeStr != "" {
		paidAmountStr = fmt.Sprintf("%s(%s)", paidAmountStr, payTypeStr)
	}

	// 优惠为空时不显示
	discountStr := ""
	if !utils.FloatEqual(discount, 0) {
		discountStr = utils.FormatFloat(discount)
	}

	// 标题取配置，默认客户名+还款凭证
	title := logic.runtime.Viper.GetString("share.title")
	if title == "" {
		title = fmt.Sprintf("%s还款凭证", customer.Name)
	}

	receiptData := utils.PaymentData{
		Title:           title,
		CustomerName:    customer.Name,
		Date:            pay.CreatedAt.Format("2006-01-02 15:04"),
		Contact:         logic.runtime.Viper.GetString("share.contact"),
		PreCreditAmount: utils.FormatFloat(preCredit),
		PaymentAmount:   paidAmountStr,
		Discount:        discountStr,
		CurrentAmount:   utils.FormatFloat(currentAmount),
	}

	html, err = utils.GeneratePaymentHTML(receiptData)
	if err != nil {
		logic.runtime.Logger.Error("GenerateShareHTML failed", zap.Error(err))
	}
	return
}
