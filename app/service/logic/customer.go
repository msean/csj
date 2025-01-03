package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (logic *CustomerLogic) Check(customer request.CustomerParam) (duplicate bool, err error) {
	var _c model.Customer
	if err = model.Find(logic.runtime.DB, &_c, model.WhereNameCond(customer.Name), model.WhereOwnerUserCond(logic.OwnerUser)); err != nil {
		return
	}
	if _c.UID != "" && _c.UID != customer.UID {
		duplicate = true
	}
	return
}

func (logic *CustomerLogic) Create(param request.CustomerParam) (err error) {
	_c := model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		Debt:      param.Debt,
		CarNo:     param.CarNo,
	}
	return model.CreateObj(logic.runtime.DB, &_c)
}

func (logic *CustomerLogic) Update(param request.CustomerParam) (err error) {
	_c := model.Customer{
		OwnerUser: logic.OwnerUser,
		Name:      param.Name,
		Phone:     param.Phone,
		Debt:      param.Debt,
		CarNo:     param.CarNo,
		BaseModel: model.BaseModel{
			UID: param.UID,
		},
	}
	return _c.Update(logic.runtime.DB)
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(searchvalue string, conds ...model.Cond) (rsp []response.ListCustomerRsp, err error) {
	var _customers []model.Customer
	conds = append(conds, []model.Cond{
		model.NewWhereCond("owner_user", logic.OwnerUser),
		model.NewOrderCond("Convert(name USING gbk)"),
	}...)
	if searchvalue != "" {
		conds = append(conds, model.NewOrLikeCond(searchvalue, model.LikeTypeBetween, "name", "phone"))
	}
	if err = model.Find(logic.runtime.DB, &_customers, conds...); err != nil {
		return
	}

	bill, _ := model.BillingCondByOwnerUser(logic.runtime.DB, logic.OwnerUser, _customers)
	for _, _customer := range _customers {
		rsp = append(rsp, response.ListCustomerRsp{
			Customer:       _customer,
			LatestBillDate: bill[_customer.UID],
		})
	}
	return
}

func LoadCustomerByUUIDList(db *gorm.DB, UUIDList []string, ownerUser string) (customerM map[string]model.Customer, err error) {
	var _customers []model.Customer
	customerM = make(map[string]model.Customer)
	if err = model.Find(db, &_customers, model.NewWhereCond("owner_user", ownerUser)); err != nil {
		return
	}
	for _, customer := range _customers {
		customerM[customer.UID] = customer
	}
	return
}
