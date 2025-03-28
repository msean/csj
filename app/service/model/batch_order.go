package model

import (
	"app/utils"
	"math"
)

const (
	BatchOrderTemp     int32 = 1   // 暂存单
	BatchOrderedCredit int32 = 2   // 赊欠单
	BatchOrderFinish   int32 = 3   // 已结算
	BatchOrderCancel   int32 = 100 //作废
	BatchOrderRefund   int32 = 101 // 退款
	BatchOrderReTurn   int32 = 102 // 退货

	BatchOrderUnshare = 1 // 未分享
	BatchOrderShared  = 2 // 已分享

	PayTypeWx   = 1 // 微信
	PayTypeZFB  = 2 // 支付宝
	PayTypeBank = 3 // 银行
	PayTypeCash = 4 // 现金
)

type (
	// 一个订单对应多个订单货品
	BatchOrder struct {
		BaseModel
		BatchUUID        int64              `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		OwnerUser        int64              `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID         int64              `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		Shared           int32              `gorm:"column:shared;comment:是否分享单" json:"shared"`
		SharedTime       int32              `gorm:"column:share_time;comment:分享时间" json:"sharedTime"`
		Status           int32              `gorm:"column:status;comment:状态" json:"status"`
		TotalAmount      float64            `gorm:"column:amount;comment:金额" json:"totalAmount"`
		CreditAmount     float64            `gorm:"column:credit_amount;comment:赊欠金额" json:"creditAmount"`
		GoodsListRelated []*BatchOrderGoods `gorm:"foreignKey:BatchOrderUID;references:UID"`
	}

	BatchOrderGoods struct {
		BaseModel
		BatchOrderUID int64 `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		BatchUUID     int64 `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		// BatchGoodsUUID string  `gorm:"column:batch_goods_uuid;comment:批次货品uuid" json:"BatchGoodsUUID"`
		GoodsUUID int64   `gorm:"column:goods_uuid;comment:货品uuid" json:"goodsUUID"`
		OwnerUser int64   `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID  int64   `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		SerialNo  string  `gorm:"column:serial_no;comment:批次序号" json:"serialNo"`
		Price     float64 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
		Weight    float64 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
		Mount     int32   `gorm:"column:mount;comment:数量" json:"mount"` // 件数
		CustomerField
		GoodsField
	}
	// 订单操作记录
	BatchOrderOpr struct {
		BaseModel
		BatchOrderUID string          `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		OwnerUser     string          `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		UserUUID      string          `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		History       utils.GormArray `gorm:"column:history;comment:历史" json:"History"`
	}
	BatchOrderPay struct {
		BaseModel
		CustomerUUID   int64   `gorm:"column:customer_uuid;comment:" json:"customerUUID"`
		OwnerUser      int64   `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		BatchOrderUUID int64   `gorm:"column:batch_order_uuid;comment:批次uuid'" json:"batchOrderUUID"`
		Amount         float64 `gorm:"column:amount;type:decimal(10,2);comment:支付金额" json:"amount"`
		PayType        int32   `gorm:"column:pay_type;comment:付款方式" json:"payType"`
	}
)

func (b *BatchOrderGoods) Amount() float64 {
	// 件数量为0 则是散装
	if b.Mount == 0 {
		return b.Price * float64(b.Weight)
	}
	return b.Price * float64(b.Mount)
}

func (b *BatchOrder) SetTotalAmount() float64 {
	var t float64
	for _, batchGoods := range b.GoodsListRelated {
		t += batchGoods.Amount()
	}
	return math.Round(float64(t))
}

func (b *BatchOrder) SetCreditAmount(pay float64) {
	b.CreditAmount -= pay
}
