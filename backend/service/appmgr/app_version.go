package appmgr

import (
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
)

type AppVersionService struct{}

// CreateAppVersion 创建App版本记录（uid 由模型 BeforeCreate 自动生成）
func (s *AppVersionService) CreateAppVersion(appVersion *appmgr.AppVersion) (err error) {
	err = global.GVA_MYSQL.Create(appVersion).Error
	return err
}

// DeleteAppVersion 删除App版本记录
func (s *AppVersionService) DeleteAppVersion(uid string) (err error) {
	err = global.GVA_MYSQL.Delete(&appmgr.AppVersion{}, "uid = ?", uid).Error
	return err
}

// DeleteAppVersionByIds 批量删除App版本记录
func (s *AppVersionService) DeleteAppVersionByIds(uids []string) (err error) {
	err = global.GVA_MYSQL.Delete(&[]appmgr.AppVersion{}, "uid in ?", uids).Error
	return err
}

// UpdateAppVersion 更新App版本记录
func (s *AppVersionService) UpdateAppVersion(appVersion appmgr.AppVersion) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.AppVersion{}).Where("uid = ?", appVersion.UID).
		Select("version", "download_addr", "is_force", "description").
		Updates(appVersion).Error
	return err
}

// GetAppVersion 根据uid获取App版本记录
func (s *AppVersionService) GetAppVersion(uid string) (appVersion appmgr.AppVersion, err error) {
	err = global.GVA_MYSQL.Where("uid = ?", uid).First(&appVersion).Error
	return
}

// GetAppVersionInfoList 分页获取App版本列表
func (s *AppVersionService) GetAppVersionInfoList(info appmgrReq.AppVersionSearch) (list []appmgr.AppVersion, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.AppVersion{})
	if info.Version != "" {
		db = db.Where("version LIKE ?", "%"+info.Version+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("created_at DESC").Find(&list).Error
	return list, total, err
}
