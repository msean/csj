package logic

import (
	"fmt"

	"app/global"
	"app/service/common"
	"app/service/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BatchOrderPayLogic struct {
	context *gin.Context
	runtime *global.RunTime
	model.BatchOrderPay
}

func NewBatchOrderPayLogic(context *gin.Context) *BatchOrderPayLogic {
	logic := &BatchOrderPayLogic{
		context: context,
		runtime: global.GlobalRunTime,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *BatchOrderPayLogic) Create(tx *gorm.DB) (err error) {
	useTxOut := true
	if logic.BatchOrderUUID == "" {
		return common.BatchOrderUUIDRequireErr
	}

	if tx == nil {
		tx = logic.runtime.DB.Begin()
		useTxOut = false
	}
	if err = model.CreateObj(tx, &logic.BatchOrderPay); err != nil {
		tx.Rollback()
		return
	}

	if err = UpdateOrderPay(tx, logic.Amount, logic.PayType, logic.BatchOrderUUID, logic.context); err != nil {
		tx.Rollback()
		return
	}

	if !useTxOut {
		tx.Commit()
	}

	return
}

func (logic *BatchOrderPayLogic) Update() (err error) {
	tx := logic.runtime.DB.Begin()
	if err = logic.BatchOrderPay.Update(tx); err != nil {
		tx.Rollback()
		return
	}
	if err = UpdateOrderPay(tx, logic.Amount, logic.PayType, logic.BatchOrderUUID, logic.context); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

// 查询该批次的单次是否结算完成，若完成，则需修改单次的状态
func UpdateOrderPay(db *gorm.DB, payFee float32, payType int32, batchOrderUUID string, ctx *gin.Context) (err error) {
	var order model.BatchOrder
	if err = model.Find(db, &order, model.WhereUIDCond(batchOrderUUID)); err != nil {
		return
	}
	order.CreditAmount = order.TotalAmount - order.CreditAmount - payFee
	if common.FloatGreat(0.0, order.CreditAmount) {
		order.Status = model.BatchOrderFinish
	}
	if err = model.WhereUIDCond(order.UID).Cond(db).Updates(&model.BatchOrder{
		CreditAmount: order.CreditAmount,
		Status:       order.Status,
	}).Error; err != nil {
		return
	}
	global.GlobalRunTime.Logger.Info(fmt.Sprintf("[BatchOrderPayLogic] [UpdateOrderPay] batch_order_uuid:%s, paidTotal: %f", batchOrderUUID, order.TotalAmount))

	// var toPayTotal float32
	// var paidTotal float32
	// var batchOrderGoodsList []model.BatchOrderGoods
	// if err = model.Find(db, &batchOrderGoodsList, model.NewWhereCond("batch_order_uuid", batchOrderUUID)); err != nil {
	// 	return
	// }
	// for _, _batchOrderGoods := range batchOrderGoodsList {
	// 	toPayTotal += _batchOrderGoods.Amount()
	// }
	// global.GlobalRunTime.Logger.Info(fmt.Sprintf("[BatchOrderPayLogic] [UpdateOrderPay] batch_order_uuid:%s, toPayTotal: %f", batchOrderUUID, toPayTotal))

	// var pays []model.BatchOrderPay
	// if err = model.Find(db, &pays, model.NewWhereCond("batch_order_uuid", batchOrderUUID)); err != nil {
	// 	return

	// }
	// for _, pay := range pays {
	// 	paidTotal += pay.Amount
	// }

	// b := NewBatchOrderLogic(ctx)
	// b.UID = batchOrderUUID
	// b.Record(true, model.HistoryStepPay, model.PayFeild{PayFee: payFee, PaidFee: paidTotal, PayType: payType})
	// if common.FloatGreat(paidTotal, toPayTotal) {
	// 	b.UID = batchOrderUUID
	// 	if err = b.BatchOrder.UpdateStatus(db, model.BatchOrderFinish); err != nil {
	// 		return
	// 	}
	// 	return
	// }
	return
}

func (logic *BatchOrderPayLogic) FromUUID() (err error) {
	return model.Find(logic.runtime.DB, &logic.BatchOrderPay)
}
