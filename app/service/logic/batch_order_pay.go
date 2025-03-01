package logic

import (
	"fmt"

	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BatchOrderPayLogic struct {
	context *gin.Context
	runtime *global.RunTime
	// model.BatchOrderPay
	OwnerUser string
}

func NewBatchOrderPayLogic(context *gin.Context) *BatchOrderPayLogic {
	logic := &BatchOrderPayLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *BatchOrderPayLogic) Create(tx *gorm.DB, param request.CreateBatchOrderPayParam, toUpdateOrder bool) (batchOrderPay model.BatchOrderPay, err error) {
	useTxOut := true
	if param.BatchOrderUUID == "" {
		err = common.BatchOrderUUIDRequireErr
		return
	}

	batchOrderPay = model.BatchOrderPay{
		CustomerUUID:   param.CustomerUUID,
		OwnerUser:      logic.OwnerUser,
		BatchOrderUUID: param.BatchOrderUUID,
		Amount:         param.Amount,
		PayType:        param.PayType,
	}
	if tx == nil {
		tx = logic.runtime.DB.Begin()
		useTxOut = false
	}
	if err = utils.GormCreateObj(tx, &batchOrderPay); err != nil {
		tx.Rollback()
		return
	}

	if toUpdateOrder {
		if err = UpdateOrderPay(tx, batchOrderPay.Amount, batchOrderPay.PayType, batchOrderPay.BatchOrderUUID, logic.context); err != nil {
			tx.Rollback()
			return
		}
	}

	if !useTxOut {
		tx.Commit()
	}

	return
}

func (logic *BatchOrderPayLogic) Update(param request.UpdateBatchOrderPayParam) (storage model.BatchOrderPay, err error) {
	if storage, err = logic.FromUUID(param.BatchOrderPayUUID); err != nil {
		return
	}
	storage.Amount = param.Amount
	storage.PayType = param.PayType
	tx := logic.runtime.DB.Begin()
	if err = dao.BatchOrderPay.Update(tx, storage); err != nil {
		tx.Rollback()
		return
	}
	if err = UpdateOrderPay(tx, storage.Amount, storage.PayType, storage.BatchOrderUUID, logic.context); err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

// 查询该批次的单次是否结算完成，若完成，则需修改单次的状态
func UpdateOrderPay(db *gorm.DB, payFee float64, payType int32, batchOrderUUID string, ctx *gin.Context) (err error) {
	var order model.BatchOrder
	if err = utils.GormFind(db, &order, utils.WhereUIDCond(batchOrderUUID)); err != nil {
		return
	}

	order.CreditAmount = order.CreditAmount - payFee
	if common.FloatGreat(0.0, order.CreditAmount) {
		order.Status = model.BatchOrderFinish
	}
	if err = utils.WhereUIDCond(order.UID).Cond(db).Updates(&model.BatchOrder{
		CreditAmount: order.CreditAmount,
		Status:       order.Status,
	}).Error; err != nil {
		return
	}
	global.Global.Logger.Info(fmt.Sprintf("[BatchOrderPayLogic] [UpdateOrderPay] batch_order_uuid:%s, paidTotal: %f", batchOrderUUID, order.TotalAmount))
	return
}

func (logic *BatchOrderPayLogic) FromUUID(uuid string) (object model.BatchOrderPay, err error) {
	err = utils.GormFind(logic.runtime.DB, &object, utils.WhereUIDCond(uuid))
	return
}
