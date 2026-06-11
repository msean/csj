package dao

import (
	"app/pkg/utils"
	"app/service/model"

	"gorm.io/gorm"
)

type customerDao struct{}

func newCustomerDao() *customerDao {
	return &customerDao{}
}

func (dao *customerDao) Update(db *gorm.DB, object model.Customer) (err error) {
	return utils.WhereUIDCond(object.UID).Cond(db).Updates(&model.Customer{
		Name:   object.Name,
		Phone:  object.Phone,
		Remark: object.Remark,
		Debt:   object.Debt,
		Addr:   object.Addr,
	}).Error
}

func (dao *customerDao) NewTempCustomer(ownerUser string, db *gorm.DB) (err error) {
	c := model.Customer{
		Name:      "现金客户",
		Remark:    "现金客户",
		OwnerUser: ownerUser,
	}
	return utils.CreateObj(db, &c)
}
