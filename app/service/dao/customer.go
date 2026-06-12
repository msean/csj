package dao

import (
	"app/pkg/utils"
	"app/service/model"
	"app/service/model/request"

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

func (dao *customerDao) NewTempCustomer(db *gorm.DB, ownerUser string) (err error) {
	c := model.Customer{
		Name:      "现金客户",
		Remark:    "现金客户",
		OwnerUser: ownerUser,
	}
	return utils.CreateObj(db, &c)
}

func (dao *customerDao) List(db *gorm.DB, ownerUser string, condtions request.CustomerListReq) (customers []model.Customer, err error) {
	conds := []utils.Cond{}
	conds = append(conds, []utils.Cond{
		utils.NewWhereCond("owner_user", ownerUser),
		utils.NewOrderCond("Convert(name USING gbk)"),
	}...)
	if condtions.SearchKey != "" {
		conds = append(conds, utils.NewOrLikeCond(condtions.SearchKey, utils.LikeTypeBetween, "name", "phone"))
	}
	err = utils.Find(db, &customers, conds...)
	return
}
