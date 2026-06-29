package logic

import (
	"fmt"

	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
	"app/service/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderPayLogic struct {
	context *gin.Context
	runtime *global.RunTime
	model.BatchOrderPay
}

func NewOrderPayLogic(context *gin.Context) *OrderPayLogic {
	logic := &OrderPayLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *OrderPayLogic) Create(tx *gorm.DB, toUpdateOrder bool) (err error) {
	if tx == nil {
		tx = logic.runtime.DB
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	if err = utils.CreateObj(tx, &logic.BatchOrderPay); err != nil {
		return
	}

	if toUpdateOrder {
		if err = UpdateOrderPay(tx, logic.Amount, logic.PayType, logic.BatchOrderUUID, logic.context); err != nil {
			return
		}
	}
	return
}

func (logic *OrderPayLogic) Update() (err error) {
	return logic.runtime.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = dao.OrderDao.UpdateOrderPay(tx, logic.OwnerUser, logic.BatchOrderPay); err != nil {
			return
		}
		err = UpdateOrderPay(tx, logic.Amount, logic.PayType, logic.BatchOrderUUID, logic.context)
		return
	})
}

// 查询该批次的单次是否结算完成，若完成，则需修改单次的状态
func UpdateOrderPay(db *gorm.DB, payFee float64, payType int32, batchOrderUUID string, ctx *gin.Context) (err error) {
	var order model.BatchOrder
	if err = utils.Find(db, &order, utils.WhereUIDCond(batchOrderUUID)); err != nil {
		return
	}

	order.CreditAmount = order.CreditAmount - payFee
	if utils.FloatGreat(0.0, order.CreditAmount) {
		order.Status = common.BatchOrderFinish
	}
	if err = utils.WhereUIDCond(order.UID).Cond(db).Updates(&model.BatchOrder{
		CreditAmount: order.CreditAmount,
		Status:       order.Status,
	}).Error; err != nil {
		return
	}
	global.Global.Logger.Info(fmt.Sprintf("[OrderPayLogic] [UpdateOrderPay] batch_order_uuid:%s, paidTotal: %f", batchOrderUUID, order.TotalAmount))
	return
}

func (logic *OrderPayLogic) FromUUID() (err error) {
	return utils.Find(logic.runtime.DB, &logic.BatchOrderPay)
}
