package logic

import (
	"app/global"
	"app/pkg/utils"
	"app/service/cache"
	"app/service/common"
	"app/service/dao"
	"app/service/model"
	"app/service/model/request"
	"app/service/model/response"
	"strings"

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
	logic.Name = strings.TrimSpace(logic.Name)
	var count int64
	query := logic.runtime.DB.Model(&model.Customer{}).Where("owner_user = ? AND name = ?", logic.OwnerUser, logic.Name)
	if logic.UID != "" {
		query = query.Where("uid <> ?", logic.UID)
	}
	err = query.Count(&count).Error
	duplicate = count > 0
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
	// 创建客户时，用 Debt 初始化 Credit（Credit 表示客户的当前总欠款额度）
	logic.Credit = logic.Debt
	return utils.CreateObj(logic.runtime.DB, &logic.Customer)
}

func (logic *CustomerLogic) Update() (err error) {
	duplicate, err := logic.Check()
	if err != nil {
		return err
	}
	if duplicate {
		return common.CustomerDuplicateErr
	}
	return logic.runtime.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = dao.CustomerDao.Update(logic.runtime.DB, logic.Customer); err != nil {
			return
		}
		return cache.CustomerCache.InvalidateCustomerCache(logic.UID, logic.OwnerUser)
	})
}

func (logic *CustomerLogic) ListCustomersByOwnerUser(conditons request.CustomerListReq) (rsp response.CustomerListRsp, err error) {
	var _customerModel []model.Customer
	if _customerModel, err = dao.CustomerDao.List(logic.runtime.DB, logic.OwnerUser, conditons); err != nil {
		return
	}
	bill, _ := dao.OrderDao.LatestOrderByCustomers(logic.runtime.DB, logic.OwnerUser, _customerModel)
	for _, _customer := range _customerModel {
		rsp.Customers = append(rsp.Customers, response.CustomerListItem{
			Customer:       _customer,
			LatestBillDate: bill[_customer.UID],
		})
	}

	// 统计商户下客户总数与总赊欠金额
	if err = logic.runtime.DB.Model(&model.Customer{}).
		Where("owner_user = ?", logic.OwnerUser).
		Count(&rsp.TotalCustomer).Error; err != nil {
		return
	}
	stats, err := dao.CustomerDao.SumOwnerStats(logic.runtime.DB, logic.OwnerUser)
	if err != nil {
		return
	}
	rsp.TotalCredit = stats.TotalDebt
	return
}

func (logic *CustomerLogic) Detail() (rsp response.CustomerDetailRsp, err error) {
	if err = global.Global.DB.Where("uid = ? AND owner_user = ? AND deleted_at IS NULL", logic.UID, logic.OwnerUser).First(&logic.Customer).Error; err != nil {
		return
	}

	rsp.Name = logic.Customer.Name
	count, amount, summaryErr := dao.CustomerPurchaseSummary(global.Global.DB, logic.OwnerUser, logic.UID, 0, 0, false)
	if summaryErr != nil {
		err = summaryErr
		return
	}
	rsp.OrderCount = float64(count)
	rsp.OrderAmount = amount
	rsp.Credit = logic.Customer.Credit

	return
}
