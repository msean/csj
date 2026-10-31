package appmgr

import (
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
)

type AppUserService struct{}

// GetAppUser 根据uid获取用户
func (s *AppUserService) GetAppUser(uid string) (user appmgr.AppUser, err error) {
	err = global.GVA_MYSQL.Where("uid = ? AND deleted_at IS NULL", uid).First(&user).Error
	return
}

// GetAppUserInfoList 分页获取用户列表
func (s *AppUserService) GetAppUserInfoList(info appmgrReq.AppUserSearch) (list []appmgr.AppUser, total int64, err error) {
	info.Page, info.PageSize = appmgrReq.NormalizePage(info.Page, info.PageSize)
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.AppUser{}).Where("deleted_at IS NULL")
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("created_at DESC, uid DESC").Find(&list).Error
	return list, total, err
}
