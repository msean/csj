package model

import (
	"app/utils"
	"time"
)

const (
	HistoryStepOrder      = 1 // 开码单
	HistoryStepOrderFix   = 2 // 修改码单
	HistoryStepOrderShare = 3 // 取码单
	HistoryStepCredit     = 4 // 赊欠
	HistoryStepPay        = 5 // 还款
	HistoryStepCash       = 6 // 收银 刚下订单的时候
	HistoryStepCrash      = 7 // 作废 退款和作废
)

type (
	BatchOrderHistory struct {
		BaseModel
		BatchOrderUID int64           `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		History       utils.GormArray `gorm:"column:history;comment:操作记录" json:"steps"`
	}
	Step struct {
		OprTime  time.Time `json:"operate_time"`
		StepType int32     `json:"step_type"`
		CustomerField
		OrderAmount  string `json:"orderAmount"`
		CreditAmount string `json:"creditAmount"`
		PayField
		GoodsList []*StepGoods `json:"goods_list"`
	}
	StepGoods struct {
		Price  float64 `json:"price"`
		Weight float64 `json:"weight"`
		Mount  int32   `json:"mount"`
		GoodsField
		Amount string `json:"amount"` //货款
	}
)
