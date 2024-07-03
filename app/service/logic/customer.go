package logic

import (
	"app/global"
	"app/service/common"
	"app/service/model"

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
		runtime: global.GlobalRunTime,
	}
	logic.OwnerUser = common.GetUserUUID(context)
	return logic
}

func (logic *CustomerLogic) Check() (duplicate bool, err error) {
	var _c model.Customer
	if err = model.Find(logic.runtime.DB, &_c, model.WhereNameCond(logic.Name), model.WhereUIDCond(logic.UID)); err != nil {
		return
	}
	if _c.UID != "" && _c.UID != logic.UID {
		duplicate = true
	}
	return
}

func (logic *CustomerLogic) Create() (err error) {
	return model.CreateObj(logic.runtime.DB, &logic.Customer)
}

func (logic *CustomerLogic) Update() (err error) {
	return logic.Customer.Update(logic.runtime.DB)
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(searchvalue string, conds ...model.Cond) (customers []CustomerLogic, err error) {
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
		customers = append(customers, CustomerLogic{
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
