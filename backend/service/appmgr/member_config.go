package appmgr

import (
	"fmt"
	"time"

	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
)

type MemberConfigService struct{}

// CreateMemberConfig 创建会员收费配置（uid 由模型 BeforeCreate 自动生成）
func (s *MemberConfigService) CreateMemberConfig(memberConfig *appmgr.MemberConfig) (err error) {
	if memberConfig.Typ < 1 || memberConfig.Typ > 3 || memberConfig.Price < 0 {
		return fmt.Errorf("会员类型或价格无效")
	}
	memberConfig.UID = ""
	err = global.GVA_MYSQL.Create(memberConfig).Error
	return err
}

// DeleteMemberConfig 删除会员收费配置
func (s *MemberConfigService) DeleteMemberConfig(uid string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.MemberConfig{}).Where("uid = ? AND deleted_at IS NULL", uid).Update("deleted_at", time.Now()).Error
	return err
}

// DeleteMemberConfigByIds 批量删除会员收费配置
func (s *MemberConfigService) DeleteMemberConfigByIds(uids []string) (err error) {
	err = global.GVA_MYSQL.Model(&appmgr.MemberConfig{}).Where("uid IN ? AND deleted_at IS NULL", uids).Update("deleted_at", time.Now()).Error
	return err
}

// UpdateMemberConfig 更新会员收费配置
func (s *MemberConfigService) UpdateMemberConfig(memberConfig appmgr.MemberConfig) (err error) {
	if memberConfig.UID == "" {
		return fmt.Errorf("uid不能为空")
	}
	if memberConfig.Typ < 1 || memberConfig.Typ > 3 || memberConfig.Price < 0 {
		return fmt.Errorf("会员类型或价格无效")
	}
	err = global.GVA_MYSQL.Model(&appmgr.MemberConfig{}).Where("uid = ? AND deleted_at IS NULL", memberConfig.UID).
		Select("type", "price").
		Updates(memberConfig).Error
	return err
}

// GetMemberConfig 根据uid获取会员收费配置
func (s *MemberConfigService) GetMemberConfig(uid string) (memberConfig appmgr.MemberConfig, err error) {
	err = global.GVA_MYSQL.Where("uid = ? AND deleted_at IS NULL", uid).First(&memberConfig).Error
	return
}

// GetMemberConfigInfoList 分页获取会员收费配置列表
func (s *MemberConfigService) GetMemberConfigInfoList(info appmgrReq.MemberConfigSearch) (list []appmgr.MemberConfig, total int64, err error) {
	info.Page, info.PageSize = appmgrReq.NormalizePage(info.Page, info.PageSize)
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.MemberConfig{}).Where("deleted_at IS NULL")
	if info.Typ != 0 {
		db = db.Where("type = ?", info.Typ)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("type ASC").Find(&list).Error
	return list, total, err
}
