package logic

import (
	"app/pkg/ticket"
	"app/pkg/utils"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func saveOrderTicketSettlement(tx *gorm.DB, order model.BatchOrder, received float64, payType int32) error {
	if order.AutoID == 0 {
		var err error
		order, err = dao.OrderDao.FindByUID(tx, order.OwnerUser, order.UID)
		if err != nil {
			return err
		}
	}
	previous, err := ticketPreviousCredit(tx, order)
	if err != nil {
		return err
	}
	settlement := model.OrderTicketSettlement{
		OwnerUser: order.OwnerUser, OrderUUID: order.UID,
		PreviousCredit: previous, OrderCredit: max(order.CreditAmount, 0),
		Received: received, PayType: payType,
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "owner_user"}, {Name: "order_uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"previous_credit", "order_credit", "received", "pay_type"}),
	}).Create(&settlement).Error
}

func ticketPreviousCredit(db *gorm.DB, order model.BatchOrder) (float64, error) {
	if order.UserUUID == "" {
		return 0, nil
	}
	customer, err := dao.CustomerDao.FromUUID(db, order.UserUUID)
	if err != nil {
		return 0, err
	}
	if customer.OwnerUser != order.OwnerUser {
		return 0, fmt.Errorf("客户不存在或无权操作")
	}
	previous, err := dao.OrderDao.HistoricalRepaymentCredit(db, order)
	if err != nil {
		return 0, err
	}
	return previous + max(customer.Debt-customer.PaidDebt, 0), nil
}

func (logic *OrderLogic) PrintTicket(req request.OrderTicketReq) (response.OrderTicketRsp, error) {
	var rsp response.OrderTicketRsp
	if req.FeedLines == nil {
		return rsp, fmt.Errorf("缺少走纸行数")
	}
	var receipt ticket.Receipt
	err := logic.runtime.DB.Transaction(func(tx *gorm.DB) error {
		var queryErr error
		receipt, queryErr = logic.ticketData(tx, req.OrderUUID)
		return queryErr
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	if err != nil {
		return rsp, err
	}
	data, err := ticket.Render(receipt, ticket.Options{
		PaperWidth: req.PaperWidth, PrintMode: req.PrintMode, FeedLines: *req.FeedLines,
	}, logic.runtime.Viper.GetString("print.font"))
	if err != nil {
		logic.runtime.Logger.Warn("小票渲染失败", zap.String("orderUUID", req.OrderUUID), zap.Error(err))
		return rsp, nil
	}
	return response.OrderTicketRsp{
		TicketBytes: base64.StdEncoding.EncodeToString(data),
		BytesLength: len(data),
	}, nil
}

func (logic *OrderLogic) ticketData(db *gorm.DB, orderUUID string) (ticket.Receipt, error) {
	var receipt ticket.Receipt
	order, err := dao.OrderDao.FindByUID(
		db.Where("owner_user = ? AND deleted_at IS NULL", logic.OwnerUser),
		logic.OwnerUser, orderUUID, "GoodsListRelated")
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return receipt, fmt.Errorf("订单不存在或无权操作")
	}
	if err != nil {
		return receipt, err
	}
	switch order.Status {
	case common.BatchOrderTemp, common.BatchOrderedCredit, common.BatchOrderFinish:
	default:
		return receipt, fmt.Errorf("当前订单状态不可打印销售小票")
	}

	owner, err := dao.UserDao.FromUUID(db, logic.OwnerUser)
	if err != nil {
		return receipt, err
	}
	if province := strings.TrimSpace(owner.Province); province != "" {
		receipt.Title = province + "食用农产品批发市场销售票据"
	}
	receipt.Market, receipt.Stall = owner.MarketName, owner.StallName
	receipt.TricycleFee, receipt.Discount = order.TricycleFee, order.Discount
	receipt.Contact = ticketContactPhone(owner.ContactPhones, owner.Phone)
	var originalCredit float64
	if order.UserUUID != "" {
		customer, err := dao.CustomerDao.FromUUID(db, order.UserUUID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return receipt, err
		}
		if customer.OwnerUser == logic.OwnerUser {
			receipt.Customer = customer.Name
			originalCredit = max(customer.Debt-customer.PaidDebt, 0)
		}
	}

	receipt.Serial = utils.UUIDToSerial(order.UID)
	orderTime := order.CreatedAt
	if order.OrderTime > 0 {
		orderTime = time.Unix(int64(order.OrderTime), 0)
	}
	receipt.Time = orderTime.In(time.FixedZone("CST", 8*60*60)).Format("2006-01-02 15:04")

	if len(order.GoodsListRelated) > 0 {
		goodsUUIDs := make([]string, 0, len(order.GoodsListRelated))
		for _, item := range order.GoodsListRelated {
			if item.DeletedAt == nil {
				goodsUUIDs = append(goodsUUIDs, item.GoodsUUID)
			}
		}
		if len(goodsUUIDs) > 0 {
			names, err := cache.GoodsCache.BatchGoodsFeildSet(goodsUUIDs, logic.OwnerUser)
			if err != nil {
				return receipt, err
			}
			for _, item := range order.GoodsListRelated {
				if item.DeletedAt != nil {
					continue
				}
				goods := names[item.GoodsUUID]
				goodsType := item.GoodType
				if goodsType == 0 {
					goodsType = goods.GoodsTyp
				}
				quantity, weight := item.Mount, item.Weight
				switch goodsType {
				case common.GoodsTypeBulk:
					quantity = 0
				case common.GoodsTypeFix:
					weight = 0
				}
				receipt.Items = append(receipt.Items, ticket.Item{
					Name: goods.GoodsName, Quantity: quantity,
					Weight: weight, Price: item.Price, Amount: item.Total,
				})
				receipt.GoodsTotal += item.Total
			}
		}
	}

	var settlement model.OrderTicketSettlement
	var orderCredit float64
	err = db.Where("owner_user = ? AND order_uuid = ?", logic.OwnerUser, order.UID).Take(&settlement).Error
	switch {
	case err == nil:
		receipt.PreviousCredit = settlement.PreviousCredit
		orderCredit = settlement.OrderCredit
		receipt.Received = settlement.Received
		if receipt.Received > 0 {
			receipt.PaymentMethod = ticketPaymentMethod(settlement.PayType)
		}
	case errors.Is(err, gorm.ErrRecordNotFound):
		orderCredit = max(order.CreditAmount, 0)
		if order.UserUUID != "" {
			if receipt.PreviousCredit, err = dao.OrderDao.HistoricalRepaymentCredit(db, order); err != nil {
				return receipt, err
			}
			receipt.PreviousCredit += originalCredit
		}
		// 旧版改单留有无还款明细的收款记录。后续还款不能当作收银实收。
		var pay model.BatchOrderPay
		err = db.Table(dao.ShardPayTable(logic.OwnerUser)).
			Where("owner_user = ? AND batch_order_uuid = ? AND customer_uuid = ''", logic.OwnerUser, order.UID).
			Where("is_revoked = 0 AND deleted_at IS NULL AND amount > 0 AND COALESCE(pay_details, '') = ''").
			Order("created_at DESC, uid DESC").Take(&pay).Error
		if err == nil {
			receipt.Received = pay.Amount
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return receipt, err
		}
	default:
		return receipt, err
	}

	receipt.TotalCredit = utils.FloatReserve(receipt.PreviousCredit+orderCredit, 2)
	return receipt, nil
}

var ticketMobile = regexp.MustCompile(`(?:^|[^0-9])([0-9]{11})(?:$|[^0-9])`)

func ticketContactPhone(contacts, accountPhone string) string {
	if strings.TrimSpace(strings.ReplaceAll(contacts, ",", "")) == "" {
		contacts = accountPhone
	}
	match := ticketMobile.FindStringSubmatch(contacts)
	if len(match) == 0 {
		return ""
	}
	return match[1]
}

// 收银客户端使用 0=微信、1=支付宝、2=现金。
func ticketPaymentMethod(payType int32) string {
	switch payType {
	case 0:
		return "微信"
	case 1:
		return "支付宝"
	case 2:
		return "现金"
	default:
		return ""
	}
}
