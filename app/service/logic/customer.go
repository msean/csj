package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/common"
	"app/service/dao"
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
	return utils.CreateObj(logic.runtime.DB, &logic.Customer)
}

func (logic *CustomerLogic) Update() (err error) {
	return dao.CustomerDao.Update(logic.runtime.DB, logic.Customer)
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(searchvalue string, conds ...utils.Cond) (customers []CustomerLogic, err error) {
	var _customers []model.Customer
	conds = append(conds, []utils.Cond{
		utils.NewWhereCond("owner_user", logic.OwnerUser),
		utils.NewOrderCond("Convert(name USING gbk)"),
	}...)
	if searchvalue != "" {
		conds = append(conds, utils.NewOrLikeCond(searchvalue, utils.LikeTypeBetween, "name", "phone"))
	}
	if err = utils.Find(logic.runtime.DB, &_customers, conds...); err != nil {
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
	if err = utils.Find(db, &_customers, utils.NewWhereCond("owner_user", ownerUser)); err != nil {
		return
	}
	for _, customer := range _customers {
		customerM[customer.UID] = customer
	}
	return
}
