package csj_customers

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/csj_customers"
	csj_customersReq "github.com/flipped-aurora/gin-vue-admin/server/model/csj_customers/request"
)

type CustomersService struct{}

// CreateCustomers 创建customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) CreateCustomers(customers *csj_customers.Customers) (err error) {
	err = global.GVA_DB.Create(customers).Error
	return err
}

// DeleteCustomers 删除customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) DeleteCustomers(ID string) (err error) {
	err = global.GVA_DB.Delete(&csj_customers.Customers{}, "uid = ?", ID).Error
	return err
}

// DeleteCustomersByIds 批量删除customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) DeleteCustomersByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]csj_customers.Customers{}, "uid in ?", IDs).Error
	return err
}

// UpdateCustomers 更新customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) UpdateCustomers(customers csj_customers.Customers) (err error) {
	err = global.GVA_DB.Model(&csj_customers.Customers{}).Where("uid = ?", customers.Uid).Updates(&customers).Error
	return err
}

// GetCustomers 根据ID获取customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) GetCustomers(ID string) (customers csj_customers.Customers, err error) {
	err = global.GVA_DB.Where("uid = ?", ID).First(&customers).Error
	return
}

// GetCustomersInfoList 分页获取customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) GetCustomersInfoList(info csj_customersReq.CustomersSearch) (list []csj_customers.Customers, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&csj_customers.Customers{})
	var customerss []csj_customers.Customers
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Value != "" {
		db = db.Where("owner_user = ?", info.Value)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&customerss).Error
	return customerss, total, err
}
