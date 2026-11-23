package logic

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const preparedTicketTTL = 30 * time.Minute

func orderTicketNo(order model.BatchOrder) string {
	if order.TicketNo != "" {
		return order.TicketNo
	}
	return utils.UUIDToSerial(order.UID)
}

func (logic *OrderLogic) preparedTicketKey() string {
	return fmt.Sprintf("order:print:%s:%s", logic.OwnerUser, logic.UserUUID)
}

func (logic *OrderLogic) printedGoodsFingerprint() string {
	items := make([]string, 0, len(logic.GoodsListRelated))
	for _, goods := range logic.GoodsListRelated {
		items = append(items, fmt.Sprintf("%q|%d|%g|%d|%g",
			goods.GoodsUUID, goods.GoodType, goods.Price, goods.Mount, goods.Weight))
	}
	sort.Strings(items)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%q|%s", logic.BatchUUID, strings.Join(items, ";")))))
}

type PreparedOrderTicket struct {
	OrderUUID string `json:"orderUUID"`
	TicketNo  string `json:"ticketNo"`
}

// Reuse the reservation for identical goods, including concurrent prepare requests.
var prepareTicketScript = redis.NewScript(`
local value = redis.call('GET', KEYS[1])
local prefix = ARGV[1] .. '*'
if value and string.sub(value, 1, string.len(prefix)) == prefix then
    redis.call('EXPIRE', KEYS[1], ARGV[3])
    return string.sub(value, string.len(prefix) + 1)
end
redis.call('SET', KEYS[1], prefix .. ARGV[2], 'EX', ARGV[3])
return ARGV[2]
`)

func (logic *OrderLogic) PrepareTicket(req request.OrderReq) (PreparedOrderTicket, error) {
	var rsp PreparedOrderTicket
	req.OwnerUser = logic.OwnerUser
	logic.BatchOrder = req.BatchOrder

	if logic.BatchUUID == "" {
		return rsp, common.BatchUUIDRequireErr
	}
	if logic.UserUUID == "" {
		return rsp, errors.New("缺少客户 UUID")
	}
	customer, err := dao.CustomerDao.FromUUID(logic.runtime.DB, logic.UserUUID)
	if err != nil {
		return rsp, err
	}
	if customer.OwnerUser != logic.OwnerUser {
		return rsp, errors.New("客户不存在或无权操作")
	}
	if err := logic.FillGoodsWeightWithTotal(logic.runtime.DB); err != nil {
		return rsp, err
	}
	if len(logic.GoodsListRelated) == 0 {
		return rsp, common.BatchOrderGoodsRequireErr
	}

	client := logic.runtime.Redis
	if client == nil {
		return rsp, errors.New("预生成票号服务不可用")
	}
	var order model.BatchOrder
	if err := order.BaseModel.BeforeCreate(logic.runtime.DB); err != nil {
		return rsp, err
	}
	uid, err := prepareTicketScript.Run(logic.context.Request.Context(), client, []string{logic.preparedTicketKey()},
		logic.printedGoodsFingerprint(), order.UID, int(preparedTicketTTL.Seconds())).Text()
	if err != nil {
		return rsp, fmt.Errorf("预生成票号失败: %w", err)
	}
	return PreparedOrderTicket{OrderUUID: uid, TicketNo: utils.UUIDToSerial(uid)}, nil
}

func (logic *OrderLogic) loadPreparedTicket() (string, error) {
	client := logic.runtime.Redis
	if client == nil {
		return "", errors.New("预生成票号服务不可用")
	}
	value, err := client.Get(logic.context.Request.Context(), logic.preparedTicketKey()).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("预生成票号不存在或已过期，请重新生成并打印")
	}
	if err != nil {
		return "", fmt.Errorf("读取预生成票号失败: %w", err)
	}
	fingerprint, uid, ok := strings.Cut(value, "*")
	if !ok || uid == "" || fingerprint != logic.printedGoodsFingerprint() {
		return "", errors.New("下单货品与预打印货品不一致，请重新生成并打印")
	}
	if logic.UID != "" && logic.UID != uid {
		return "", errors.New("预生成订单 UUID 已变化，请重新生成并打印")
	}
	logic.UID = uid
	return value, nil
}

var consumePreparedTicketScript = redis.NewScript(`
if redis.call('GET', KEYS[1]) == ARGV[1] then
    return redis.call('DEL', KEYS[1])
end
return 0
`)

func (logic *OrderLogic) consumePreparedTicket(value string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := consumePreparedTicketScript.Run(ctx, logic.runtime.Redis, []string{logic.preparedTicketKey()}, value).Err(); err != nil {
		logic.runtime.Logger.Warn("清理预生成票号失败", zap.String("orderUUID", logic.UID), zap.Error(err))
	}
}
