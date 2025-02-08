package model

import (
	"app/service/common"
	"app/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
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
		BatchOrderUID string          `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
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

func (bo *BatchOrder) Record(db *gorm.DB, stepType int32, pay PayField) (err error) {
	var boh BatchOrderHistory
	if err = utils.GormFind(db, &boh, utils.NewWhereCond("batch_order_uuid", bo.UID)); err != nil {
		return
	}
	step := bo.NewHistoryStep(stepType, pay)
	boh.History = append(boh.History, step)
	boh.BatchOrderUID = bo.UID
	if boh.UID == "" {
		return utils.GormCreateObj(db, &boh)
	}
	return db.Save(&boh).Error
}

func (bo *BatchOrder) NewHistoryStep(stepType int32, pay PayField) (step Step) {
	s := Step{
		CustomerField: bo.CustomerField,
		StepType:      stepType,
		OprTime:       time.Now(),
	}
	for _, goods := range bo.GoodsListRelated {
		s.GoodsList = append(s.GoodsList, &StepGoods{
			Price:      goods.Price,
			Weight:     goods.Weight,
			Mount:      goods.Mount,
			GoodsField: goods.GoodsField,
			Amount:     common.Float32Preserve(goods.Amount(), 2),
		})
	}
	s.CreditAmount = fmt.Sprintf("%.f", bo.CreditAmount)
	s.PayField = pay
	// if !common.Float32IsZero(pay.PaidFee) && pay.PaidFee > 0.0 {
	// 	s.CreditAmount = common.Float32Preserve((_orderAmount - pay.PaidFee), 32)
	// }
	step = s
	return
}
