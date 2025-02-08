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
)

type CustomerLogic struct {
	context   *gin.Context
	runtime   *global.RunTime
	OwnerUser string
}

func NewCustomerLogic(context *gin.Context) *CustomerLogic {
	logic := &CustomerLogic{
		context: context,
		runtime: global.Global,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *CustomerLogic) Check(param request.CustomerParam) (duplicate bool, err error) {
	var _c model.Customer
	if err = utils.GormFind(logic.runtime.DB, &_c, utils.WhereNameCond(param.Name), utils.WhereUIDCond(param.UID)); err != nil {
		return
	}
	if _c.UID != "" && _c.UID != param.UID {
		duplicate = true
	}
	return
}

func (logic *CustomerLogic) Create(param request.CustomerParam) (err error) {
	customerModel := model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		CarNo:     param.CarNo,
		Debt:      param.Debt,
	}
	return utils.GormCreateObj(logic.runtime.DB, &customerModel)
}

func (logic *CustomerLogic) Update(param request.CustomerParam) (err error) {
	customerModel := model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		CarNo:     param.CarNo,
		Debt:      param.Debt,
	}
	customerModel.UID = param.UID
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
