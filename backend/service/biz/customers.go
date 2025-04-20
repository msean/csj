package biz

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/flipped-aurora/gin-vue-admin/server/dao"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
	biz_request "github.com/flipped-aurora/gin-vue-admin/server/model/biz/request"
)

type CustomersService struct{}

// CreateCustomers 创建customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) CreateCustomers(customers *biz.Customers) (err error) {
	err = global.GVA_DB.Create(customers).Error
	return err
}

// DeleteCustomers 删除customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) DeleteCustomers(ID string) (err error) {
	err = global.GVA_DB.Delete(&biz.Customers{}, "uid = ?", ID).Error
	return err
}

// DeleteCustomersByIds 批量删除customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) DeleteCustomersByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]biz.Customers{}, "uid in ?", IDs).Error
	return err
}

// UpdateCustomers 更新customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) UpdateCustomers(customers biz.Customers) (err error) {
	err = global.GVA_DB.Model(&biz.Customers{}).Where("uid = ?", customers.Uid).Updates(&customers).Error
	return err
}

// GetCustomers 根据ID获取customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) GetCustomers(ID string) (customers biz.Customers, err error) {
	err = global.GVA_DB.Where("uid = ?", ID).First(&customers).Error
	return
}

// GetCustomersInfoList 分页获取customers表记录
// Author [piexlmax](https://github.com/piexlmax)
func (customersService *CustomersService) GetCustomersInfoList(info biz_request.CustomersSearch) (list []*biz.Customers, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&biz.Customers{})
	var customerss []*biz.Customers
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
	var ownerUsers []int64
	for _, customer := range customerss {
		ownerUsers = append(ownerUsers, customer.OwnerUserUUID)
	}
	userMapper := make(map[int64]biz.Users)
	if userMapper, err = dao.MapperByOwnerUser(ownerUsers); err != nil {
		return
	}

	spew.Dump(">>>>>>>>>", userMapper)
	for _, customer := range customerss {
		customer.OwnerUser = userMapper[customer.OwnerUserUUID].Name
	}
	return customerss, total, err
}
