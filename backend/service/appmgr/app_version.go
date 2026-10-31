package appmgr

import (
	"fmt"
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
	"strings"
	"time"
	"unicode/utf8"
)

type AppVersionService struct{}

// CreateAppVersion 创建App版本记录（uid 由模型 BeforeCreate 自动生成）
func (s *AppVersionService) CreateAppVersion(appVersion *appmgr.AppVersion) (err error) {
	if strings.TrimSpace(appVersion.Version) == "" || strings.TrimSpace(appVersion.DownloadAddr) == "" || utf8.RuneCountInString(appVersion.Version) > 100 || utf8.RuneCountInString(appVersion.DownloadAddr) > 100 || utf8.RuneCountInString(appVersion.Description) > 500 {
		return fmt.Errorf("版本号、下载地址不能为空且不能超长")
	}
	appVersion.UID = ""
	err = global.GVA_MYSQL.Create(appVersion).Error
	return err
}

// DeleteAppVersion 删除App版本记录
func (s *AppVersionService) DeleteAppVersion(uid string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.AppVersion{}).Where("uid = ? AND deleted_at IS NULL", uid).Update("deleted_at", time.Now()).Error
	return err
}

// DeleteAppVersionByIds 批量删除App版本记录
func (s *AppVersionService) DeleteAppVersionByIds(uids []string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.AppVersion{}).Where("uid IN ? AND deleted_at IS NULL", uids).Update("deleted_at", time.Now()).Error
	return err
}

// UpdateAppVersion 更新App版本记录
func (s *AppVersionService) UpdateAppVersion(appVersion appmgr.AppVersion) (err error) {
	if appVersion.UID == "" {
		return fmt.Errorf("uid不能为空")
	}
	if strings.TrimSpace(appVersion.Version) == "" || strings.TrimSpace(appVersion.DownloadAddr) == "" || utf8.RuneCountInString(appVersion.Version) > 100 || utf8.RuneCountInString(appVersion.DownloadAddr) > 100 || utf8.RuneCountInString(appVersion.Description) > 500 {
		return fmt.Errorf("版本号、下载地址不能为空且不能超长")
	}
	err = global.GVA_MYSQL.Model(&appmgr.AppVersion{}).Where("uid = ? AND deleted_at IS NULL", appVersion.UID).
		Select("version", "download_addr", "is_force", "description").
		Updates(appVersion).Error
	return err
}

// GetAppVersion 根据uid获取App版本记录
func (s *AppVersionService) GetAppVersion(uid string) (appVersion appmgr.AppVersion, err error) {
	err = global.GVA_MYSQL.Where("uid = ? AND deleted_at IS NULL", uid).First(&appVersion).Error
	return
}

// GetAppVersionInfoList 分页获取App版本列表
func (s *AppVersionService) GetAppVersionInfoList(info appmgrReq.AppVersionSearch) (list []appmgr.AppVersion, total int64, err error) {
	info.Page, info.PageSize = appmgrReq.NormalizePage(info.Page, info.PageSize)
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.AppVersion{}).Where("deleted_at IS NULL")
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
	err = db.Order("created_at DESC, uid DESC").Find(&list).Error
	return list, total, err
}
