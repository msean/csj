package dao

import (
	"app/service/common"
	"app/service/model"
	"app/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BatchOrderDao struct{}

func NewBatchOrderDao() *BatchOrderDao {
	return &BatchOrderDao{}
}

func (dao *BatchOrderDao) UpdateStatus(db *gorm.DB, uid int64, status int32) error {
	return utils.WhereUIDCond(uid).Cond(db).Model(&model.BatchOrder{}).Update("status", status).Error
}

func (dao *BatchOrderDao) UpdateShare(db *gorm.DB, uid int64, share int32) error {
	return utils.WhereUIDCond(uid).Cond(db).Updates(&model.BatchOrder{
		Shared: model.BatchOrderShared,
	}).Error
}

func (bo *BatchOrderDao) Record(
	db *gorm.DB,
	batchOrder model.BatchOrder,
	batchOrderGoodsList []model.BatchOrderGoods,
	stepType int32,
	pay model.PayField,
	customerField model.CustomerField,
) (err error) {
	var boh model.BatchOrderHistory
	if err = utils.GormFind(db, &boh, utils.NewWhereCond("batch_order_uuid", batchOrder.UID)); err != nil {
		return
	}
	step := bo.NewHistoryStep(batchOrder, stepType, pay, batchOrderGoodsList, customerField)
	boh.History = append(boh.History, step)
	boh.BatchOrderUID = batchOrder.UID
	if boh.UID == 0 {
		return utils.GormCreateObj(db, &boh)
	}
	return db.Save(&boh).Error
}

func (bo *BatchOrderDao) NewHistoryStep(batchOrder model.BatchOrder, stepType int32, pay model.PayField, batchOrderGoodsList []model.BatchOrderGoods, customerField model.CustomerField) (step model.Step) {
	s := model.Step{
		CustomerField: customerField,
		StepType:      stepType,
		OprTime:       time.Now(),
	}

	for _, goods := range batchOrderGoodsList {
		s.GoodsList = append(s.GoodsList, &model.StepGoods{
			Price:      goods.Price,
			Weight:     goods.Weight,
			Mount:      goods.Mount,
			GoodsField: goods.GoodsField,
			Amount:     common.Float32Preserve(goods.Amount(), 2),
		})
	}
	s.CreditAmount = fmt.Sprintf("%.f", batchOrder.CreditAmount)
	s.PayField = pay
	// if !common.Float32IsZero(pay.PaidFee) && pay.PaidFee > 0.0 {
	// 	s.CreditAmount = common.Float32Preserve((_orderAmount - pay.PaidFee), 32)
	// }
	step = s
	return
}
