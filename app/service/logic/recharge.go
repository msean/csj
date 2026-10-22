package logic

import (
	"app/global"
	"app/pkg/wechat"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/response"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RechargeLogic struct {
	context *gin.Context
	runtime *global.RunTime
}

func NewRechargeLogic(context *gin.Context) *RechargeLogic {
	return &RechargeLogic{
		context: context,
		runtime: global.Global,
	}
}

// memberTypeName 会员类型描述
func memberTypeName(memberType int) string {
	switch memberType {
	case common.MemberTypeMonth:
		return "会员充值-月卡"
	case common.MemberTypeQuarter:
		return "会员充值-季卡"
	case common.MemberTypeYear:
		return "会员充值-年卡"
	default:
		return "会员充值"
	}
}

// extendDuration 根据会员类型返回延长时长
func extendDuration(memberType int) (time.Duration, error) {
	switch memberType {
	case common.MemberTypeMonth:
		return time.Hour * 24 * 30, nil
	case common.MemberTypeQuarter:
		return time.Hour * 24 * 90, nil
	case common.MemberTypeYear:
		return time.Hour * 24 * 365, nil
	default:
		return 0, common.MemberTypeInvalidErr
	}
}

// generateOrderNo 生成商户订单号
func generateOrderNo() string {
	return fmt.Sprintf("CSJ%s%04d", time.Now().Format("20060102150405"), rand.Intn(10000))
}

// CreateRechargeOrder 创建会员充值订单
func (logic *RechargeLogic) CreateRechargeOrder(memberType int) (rsp response.CreateRechargeRsp, err error) {
	// 校验会员类型
	if _, e := extendDuration(memberType); e != nil {
		err = e
		return
	}

	// 查询会员配置价格
	var config model.MemberConfig
	if err = logic.runtime.DB.Where("type = ?", memberType).First(&config).Error; err != nil {
		err = common.MemberTypeInvalidErr
		return
	}

	userUUID := common.GetUserUUID(logic.context)
	orderNo := generateOrderNo()
	amountFen := int64(config.Price * 100)

	// 创建本地充值订单
	order := &model.RechargeOrder{
		UserUUID:   userUUID,
		OrderNo:    orderNo,
		MemberType: memberType,
		Amount:     config.Price,
		PayStatus:  0,
	}
	if err = dao.RechargeDao.Create(logic.runtime.DB, order); err != nil {
		return
	}

	// 调用微信支付下单
	payInfo, e := wechat.AppPay(logic.context, orderNo, memberTypeName(memberType), amountFen)
	if e != nil {
		// 下单失败不删除订单，保留待支付状态
		err = e
		return
	}

	rsp.OrderNo = orderNo
	rsp.PayInfo = payInfo
	return
}

// WechatNotify 处理微信支付回调
func (logic *RechargeLogic) WechatNotify() error {
	tx, err := wechat.ParseNotify(logic.context, logic.context.Request)
	if err != nil {
		return err
	}

	if tx.OutTradeNo == nil || tx.TradeState == nil {
		return fmt.Errorf("回调数据不完整")
	}

	orderNo := *tx.OutTradeNo
	tradeState := *tx.TradeState

	// 仅处理支付成功
	if tradeState != "SUCCESS" {
		return nil
	}

	// 事务：标记订单已支付 + 延长会员时间
	return logic.runtime.DB.Transaction(func(txDb *gorm.DB) error {
		order, e := dao.RechargeDao.MarkPaid(txDb, orderNo, *tx.TransactionId)
		if e != nil {
			return e
		}

		// 幂等：如果订单已支付（非本次更新），直接返回
		// MarkPaid 已保证幂等，这里只需延长会员时间
		user, e := dao.UserDao.FromUUID(txDb, order.UserUUID)
		if e != nil {
			return e
		}

		duration, e := extendDuration(order.MemberType)
		if e != nil {
			return e
		}

		// 计算新的截止时间：若当前截止时间在未来，则从截止时间起延长；否则从当前时间起延长
		baseTime := time.Now()
		if user.VipDeadline.After(baseTime) {
			baseTime = user.VipDeadline
		}
		newDeadline := baseTime.Add(duration)

		return dao.UserDao.UpdateVipDeadline(txDb, order.UserUUID, newDeadline)
	})
}

// QueryRechargeStatus 查询充值订单状态
func (logic *RechargeLogic) QueryRechargeStatus(orderNo string) (rsp response.RechargeStatusRsp, err error) {
	order, err := dao.RechargeDao.FindByOrderNo(logic.runtime.DB, orderNo)
	if err != nil {
		err = common.RechargeOrderNotExist
		return
	}
	rsp.OrderNo = order.OrderNo
	rsp.PayStatus = order.PayStatus
	rsp.PayTime = order.PayTime

	userUUID := common.GetUserUUID(logic.context)
	if userUUID != "" {
		user, e := dao.UserDao.FromUUID(logic.runtime.DB, userUUID)
		if e == nil {
			rsp.VipDeadline = user.VipDeadline
		}
	}
	return
}

// RechargeList 当前用户充值记录
func (logic *RechargeLogic) RechargeList(page, pageCount int) (list []model.RechargeOrder, total int64, err error) {
	if page <= 0 {
		page = 1
	}
	if pageCount <= 0 {
		pageCount = 20
	}
	userUUID := common.GetUserUUID(logic.context)
	return dao.RechargeDao.ListByUser(logic.runtime.DB, userUUID, page, pageCount)
}
