package dao

import (
	"app/pkg/utils"
	"app/service/common"
	"app/service/model"

	"gorm.io/gorm"
)

type orderDao struct{}

func newOrderDao() *orderDao {
	return &orderDao{}
}

func (dao *orderDao) UpdateStatus(db *gorm.DB, orderUUID string, status int) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Model(&model.BatchOrder{}).Update("status", status).Error
}

func (dao *orderDao) Shared(db *gorm.DB, orderUUID string) error {
	return utils.WhereUIDCond(orderUUID).Cond(db).Updates(&model.BatchOrder{
		Shared: common.BatchOrderShared,
	}).Error
}

func (dao *orderDao) UpdateGoods(db *gorm.DB, object model.BatchOrderGoods) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderGoods{
		Price:  object.Price,
		Mount:  object.Mount,
		Weight: object.Weight,
	}).Error
}

func (dao *orderDao) UpdateOrderPay(db *gorm.DB, object model.BatchOrderPay) error {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.BatchOrderPay{
		PayType: object.PayType,
		Amount:  object.Amount,
	}).Error
}
