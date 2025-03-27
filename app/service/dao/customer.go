package dao

import (
	"app/service/model"
	"app/utils"

	"gorm.io/gorm"
)

type CustomerDao struct{}

func NewCustomerDao() *CustomerDao {
	return &CustomerDao{}
}

func (dao *CustomerDao) Update(db *gorm.DB, customer model.Customer) (err error) {
	return utils.WhereUIDCond(customer.UID).Cond(db).Updates(&model.Customer{
		Name:   customer.Name,
		Phone:  customer.Phone,
		Remark: customer.Remark,
		Debt:   customer.Debt,
		Addr:   customer.Addr,
	}).Error
}

func (dao *CustomerDao) NewTempCustomer(ownerUser int64, db *gorm.DB) (err error) {
	customer := model.Customer{
		Name:      "临时客户",
		Remark:    "临时客户",
		OwnerUser: ownerUser,
	}
	return utils.GormCreateObj(db, &customer)
}

func (c *CustomerDao) MapperFromList(db *gorm.DB, UUIDList []int64, ownerUser int64) (customerM map[int64]model.Customer, err error) {
	var _customers []model.Customer
	customerM = make(map[int64]model.Customer)
	conditions := []utils.Cond{
		utils.NewWhereCond("owner_user", ownerUser),
		utils.NewInCondFromInt64("uid", UUIDList),
	}
	if err = utils.GormFind(db, &_customers, conditions...); err != nil {
		return
	}
	for _, customer := range _customers {
		customerM[customer.UID] = customer
	}
	return
}

func (c *CustomerDao) FromUUID(db *gorm.DB, uuid, ownerUser int64) (customer model.Customer, err error) {
	var _customers model.Customer
	conditions := []utils.Cond{
		utils.NewWhereCond("owner_user", ownerUser),
		utils.NewWhereCond("uid", uuid),
	}
	if err = utils.GormFind(db, &_customers, conditions...); err != nil {
		return
	}
	return
}
