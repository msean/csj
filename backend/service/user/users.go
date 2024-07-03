package user

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/user"
	userReq "github.com/flipped-aurora/gin-vue-admin/server/model/user/request"
)

type UsersService struct{}

// CreateUsers 创建users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) CreateUsers(users *user.Users) (err error) {
	err = global.GVA_DB.Create(users).Error
	return err
}

// DeleteUsers 删除users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) DeleteUsers(ID string) (err error) {
	err = global.GVA_DB.Delete(&user.Users{}, "uid = ?", ID).Error
	return err
}

// DeleteUsersByIds 批量删除users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) DeleteUsersByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]user.Users{}, "uid in ?", IDs).Error
	return err
}

// UpdateUsers 更新users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) UpdateUsers(users user.Users) (err error) {
	err = global.GVA_DB.Model(&user.Users{}).Where("uid = ?", users.Uid).Updates(&users).Error
	return err
}

// GetUsers 根据ID获取users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) GetUsers(UID string) (users user.Users, err error) {
	err = global.GVA_DB.Where("uid = ?", UID).First(&users).Error
	return
}

// GetUsersInfoList 分页获取users表记录
// Author [piexlmax](https://github.com/piexlmax)
func (usersService *UsersService) GetUsersInfoList(info userReq.UsersSearch) (list []user.Users, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&user.Users{})
	var userss []user.Users
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Value != "" {
		db = db.Where(fmt.Sprintf("phone like '%s%s%s' or name like '%s%s%s'", "%", info.Value, "%", "%", info.Value, "%"))
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&userss).Error
	return userss, total, err
}
