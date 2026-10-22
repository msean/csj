package appmgr

import (
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/appmgr"
	appmgrReq "github.com/msean/csj/backend/model/appmgr/request"
)

type RechargeOrderService struct{}

// GetRechargeOrder 根据uid获取充值订单
func (s *RechargeOrderService) GetRechargeOrder(uid string) (order appmgr.RechargeOrder, err error) {
	err = global.GVA_MYSQL.Where("uid = ?", uid).First(&order).Error
	return
}

// GetRechargeOrderInfoList 分页获取充值订单列表
func (s *RechargeOrderService) GetRechargeOrderInfoList(info appmgrReq.RechargeOrderSearch) (list []appmgr.RechargeOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_MYSQL.Model(&appmgr.RechargeOrder{})
	if info.UserUUID != "" {
		db = db.Where("user_uuid = ?", info.UserUUID)
	}
	if info.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+info.OrderNo+"%")
	}
	if info.PayStatus != 0 {
		db = db.Where("pay_status = ?", info.PayStatus)
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
