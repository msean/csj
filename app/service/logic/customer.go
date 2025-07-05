package logic

import (
	"app/global"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"app/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CustomerLogic struct {
	context   *gin.Context
	runtime   *global.RunTime
	OwnerUser int64
}

func NewCustomerLogic(context *gin.Context) *CustomerLogic {
	logic := &CustomerLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *CustomerLogic) ExistName(name string) (duplicate bool, err error) {
	var _c model.Customer
	if err = utils.GormFind(logic.runtime.DB, &_c, utils.WhereNameCond(name), utils.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		logic.runtime.Logger.Error("CustomerLogic Exist", zap.Int64("owner_user", logic.OwnerUser), zap.String("name", name), zap.Error(err))
		return
	}
	if _c.UID != 0 {
		duplicate = true
	}
	return
}

func (logic *CustomerLogic) Create(param request.CustomerParam) (customerModel model.Customer, err error) {
	customerModel = model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		CarNo:     param.CarNo,
		Debt:      param.Debt,
	}
	err = utils.GormCreateObj(logic.runtime.DB, &customerModel)
	return
}

func (logic *CustomerLogic) Update(param request.CustomerParam) (err error) {
	customerModel := model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		CarNo:     param.CarNo,
		Debt:      param.Debt,
	}
	customerModel.UID = param.UIDCompatible
	return dao.Customer.Update(logic.runtime.DB, customerModel)
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(searchvalue string, conds ...utils.Cond) (customers []response.ListCustomerRsp, err error) {
	var _customers []model.Customer
	conds = append(conds, []utils.Cond{
		utils.NewWhereCond("owner_user", logic.OwnerUser),
		utils.NewOrderCond("Convert(name USING gbk)"),
	}...)
	if searchvalue != "" {
		conds = append(conds, utils.NewOrLikeCond(searchvalue, utils.LikeTypeBetween, "name", "phone"))
	}
	if err = utils.GormFind(logic.runtime.DB, &_customers, conds...); err != nil {
		return
	}

	bill, _ := model.BillingCondByOwnerUser(logic.runtime.DB, logic.OwnerUser, _customers)
	for _, _customer := range _customers {
		customers = append(customers, response.ListCustomerRsp{
			Customer:       _customer,
			LatestBillDate: bill[_customer.UID],
		})
	}
	return
}

func (logic *CustomerLogic) FromUUID(uuid int64) (rsp response.ListCustomerRsp, err error) {
	var customer model.Customer
	if customer, err = dao.Customer.FromUUID(logic.runtime.DB, uuid, logic.OwnerUser); err != nil {
		return
	}

	bill, _ := model.BillingCondByOwnerUser(logic.runtime.DB, logic.OwnerUser, []model.Customer{customer})
	rsp = response.ListCustomerRsp{
		Customer:       customer,
		LatestBillDate: bill[customer.UID],
	}
	return
}
