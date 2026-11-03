package logic

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"app/service/dao"
	"app/service/model"
	"app/service/model/response"
	"gorm.io/gorm"
)

func decodeRepaymentAllocations(raw string) []response.RepaymentAllocation {
	rows := make([]response.RepaymentAllocation, 0)
	var value interface{} = raw
	for i := 0; i < 3; i++ {
		switch v := value.(type) {
		case string:
			if json.Unmarshal([]byte(v), &value) != nil {
				return rows
			}
		case map[string]interface{}:
			if v["orderUUID"] != nil || v["type"] != nil {
				value = []interface{}{v}
			} else if v["items"] != nil {
				value = v["items"]
			} else if v["details"] != nil {
				value = v["details"]
			} else {
				value = v["payDetails"]
			}
		default:
			i = 3
		}
	}
	list, ok := value.([]interface{})
	if !ok {
		return rows
	}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		amount, err := strconv.ParseFloat(fmt.Sprint(m["amount"]), 64)
		if err != nil || math.IsNaN(amount) || math.IsInf(amount, 0) {
			continue
		}
		typ, _ := strconv.Atoi(fmt.Sprint(m["type"]))
		uuid, _ := m["orderUUID"].(string)
		if typ == 0 && uuid != "" {
			typ = 2
		}
		rows = append(rows, response.RepaymentAllocation{Type: typ, OrderUUID: uuid, Amount: amount})
	}
	return rows
}

func enrichRepaymentAllocations(rows []response.RepaymentAllocation, orders []model.BatchOrder) {
	byID := make(map[string]model.BatchOrder, len(orders))
	for _, order := range orders {
		byID[order.UID] = order
	}
	for i := range rows {
		if rows[i].Type == 1 {
			continue
		}
		order, ok := byID[rows[i].OrderUUID]
		if !ok {
			continue
		}
		total := order.TotalAmount + order.TricycleFee
		when := order.CreatedAt
		if order.OrderTime > 0 {
			when = time.Unix(int64(order.OrderTime), 0)
		}
		rows[i].OrderAvailable = true
		rows[i].OrderTotalAmount = &total
		rows[i].OrderStatus = &order.Status
		if !when.IsZero() {
			rows[i].OrderTime = &when
		}
	}
}

func repaymentAllocations(db *gorm.DB, owner string, pay *model.BatchOrderPay) ([]response.RepaymentAllocation, error) {
	rows := decodeRepaymentAllocations(pay.PayDetails)
	raw := strings.TrimSpace(pay.PayDetails)
	if len(rows) == 0 && (raw == "" || raw == "null" || raw == "[]") && pay.BatchOrderUUID != "" {
		rows = append(rows, response.RepaymentAllocation{Type: 2, OrderUUID: pay.BatchOrderUUID, Amount: pay.Amount + pay.Discount})
	}
	ids := make([]string, 0)
	seen := map[string]bool{}
	for _, row := range rows {
		if row.Type != 1 && row.OrderUUID != "" && !seen[row.OrderUUID] {
			ids = append(ids, row.OrderUUID)
			seen[row.OrderUUID] = true
		}
	}
	if len(ids) == 0 {
		return rows, nil
	}
	var orders []model.BatchOrder
	// Both owner and customer are checked, even if a historical allocation has a bad link.
	err := db.Table(dao.ShardOrderTable(owner)).Where("owner_user = ? AND user_uuid = ? AND uid IN ? AND deleted_at IS NULL", owner, pay.CustomerUUID, ids).
		Select("uid, amount, tricycle_fee, order_time, created_at, status").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	enrichRepaymentAllocations(rows, orders)
	return rows, nil
}
