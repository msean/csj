package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerLogic struct {
	context *gin.Context
	runtime *global.RunTime
	model.Customer
	LatestBillDate int `json:"latestBillDate"`
}

func NewCustomerLogic(context *gin.Context) *CustomerLogic {
	logic := &CustomerLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *CustomerLogic) Check() (duplicate bool, err error) {
	var _c model.Customer
	if err = utils.Find(logic.runtime.DB, &_c, utils.WhereNameCond(logic.Name), utils.WhereUIDCond(logic.UID)); err != nil {
		return
	}
	if _c.UID != "" && _c.UID != logic.UID {
		duplicate = true
	}
	return
}

func (logic *CustomerLogic) Create() (err error) {
	var dup bool
	if dup, err = logic.Check(); err != nil {
		return
	}
	if dup {
		err = common.CustomerDuplicateErr
		return
	}
	return utils.CreateObj(logic.runtime.DB, &logic.Customer)
}

func (logic *CustomerLogic) Update() (err error) {
	return logic.runtime.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = dao.CustomerDao.Update(logic.runtime.DB, logic.Customer); err != nil {
			return
		}
		return cache.CustomerCache.InvalidateCustomerCache(logic.UID, logic.OwnerUser)
	})
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(conditons request.CustomerListReq) (customers []CustomerLogic, err error) {
	var _customerModel []model.Customer
	if _customerModel, err = dao.CustomerDao.List(logic.runtime.DB, logic.OwnerUser, conditons); err != nil {
		return
	}
	bill, _ := dao.OrderDao.LatestOrderByCustomers(logic.runtime.DB, logic.OwnerUser, _customerModel)
	for _, _customer := range _customerModel {
		customers = append(customers, CustomerLogic{
			Customer:       _customer,
			LatestBillDate: bill[_customer.UID],
		})
	}
	return
}
