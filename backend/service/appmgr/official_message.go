package appmgr

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
)

type OfficialMessageService struct{}

// CreateOfficialMessage 创建官方消息（uid 由模型 BeforeCreate 自动生成）
func (s *OfficialMessageService) CreateOfficialMessage(officialMessage *appmgr.OfficialMessage) (err error) {
	if err = validateOfficialMessage(officialMessage); err != nil {
		return
	}
	officialMessage.UID = ""
	err = global.GVA_MYSQL.Create(officialMessage).Error
	return err
}

// DeleteOfficialMessage 删除官方消息
func (s *OfficialMessageService) DeleteOfficialMessage(uid string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.OfficialMessage{}).Where("uid = ? AND deleted_at IS NULL", uid).Update("deleted_at", time.Now()).Error
	return err
}

// DeleteOfficialMessageByIds 批量删除官方消息
func (s *OfficialMessageService) DeleteOfficialMessageByIds(uids []string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.OfficialMessage{}).Where("uid IN ? AND deleted_at IS NULL", uids).Update("deleted_at", time.Now()).Error
	return err
}

// UpdateOfficialMessage 更新官方消息
func (s *OfficialMessageService) UpdateOfficialMessage(officialMessage appmgr.OfficialMessage) (err error) {
	if officialMessage.UID == "" {
		return fmt.Errorf("uid不能为空")
	}
	if err = validateOfficialMessage(&officialMessage); err != nil {
		return
	}
	err = global.GVA_MYSQL.Model(&appmgr.OfficialMessage{}).Where("uid = ? AND deleted_at IS NULL", officialMessage.UID).
		Select("title", "content", "status", "published_at", "expires_at", "sort", "level").
		Updates(officialMessage).Error
	return err
}

// GetOfficialMessage 根据uid获取官方消息
func (s *OfficialMessageService) GetOfficialMessage(uid string) (officialMessage appmgr.OfficialMessage, err error) {
	err = global.GVA_MYSQL.Where("uid = ? AND deleted_at IS NULL", uid).First(&officialMessage).Error
	return
}

// GetOfficialMessageInfoList 分页获取官方消息列表
func (s *OfficialMessageService) GetOfficialMessageInfoList(info appmgrReq.OfficialMessageSearch) (list []appmgr.OfficialMessage, total int64, err error) {
	if info.Page < 1 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	if info.PageSize > 100 {
		info.PageSize = 100
	}
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.OfficialMessage{}).Where("deleted_at IS NULL")
	if info.Status != nil {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("sort DESC, published_at DESC, uid DESC").Find(&list).Error
	return list, total, err
}

func validateOfficialMessage(m *appmgr.OfficialMessage) error {
	if strings.TrimSpace(m.Title) == "" || strings.TrimSpace(m.Content) == "" {
		return fmt.Errorf("标题和内容不能为空")
	}
	if utf8.RuneCountInString(m.Title) > 200 {
		return fmt.Errorf("标题不能超过200字")
	}
	if m.Level == 0 {
		m.Level = 1
	}
	if m.Level != 1 && m.Level != 2 {
		return fmt.Errorf("无效的通知等级")
	}
	if m.Status != 0 && m.Status != 1 {
		return fmt.Errorf("无效的发布状态")
	}
	if m.PublishedAt.IsZero() {
		m.PublishedAt = time.Now()
	}
	if m.ExpiresAt != nil && !m.ExpiresAt.After(m.PublishedAt) {
		return fmt.Errorf("失效时间必须晚于发布时间")
	}
	return nil
}
