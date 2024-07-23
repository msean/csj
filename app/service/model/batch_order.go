package model

import (
	"math"

	"gorm.io/gorm"
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
		BatchUUID        string             `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		OwnerUser        string             `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID         string             `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		Shared           int32              `gorm:"column:shared;comment:是否分享单" json:"shared"`
		SharedTime       int32              `gorm:"column:share_time;comment:分享时间" json:"sharedTime"`
		Status           int32              `gorm:"column:status;comment:状态" json:"status"`
		TotalAmount      float32            `gorm:"column:amount;comment:金额" json:"total_amount"`
		CreditAmount     float32            `gorm:"column:credit_amount;comment:赊欠金额" json:"credit_amount"`
		GoodsListRelated []*BatchOrderGoods `gorm:"foreignKey:BatchOrderUID;references:UID" json:"goodsList"`
		CustomerFeild
	}

	BatchOrderGoods struct {
		BaseModel
		BatchOrderUID string `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		BatchUUID     string `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		// BatchGoodsUUID string  `gorm:"column:batch_goods_uuid;comment:批次货品uuid" json:"BatchGoodsUUID"`
		GoodsUUID string  `gorm:"column:goods_uuid;comment:货品uuid" json:"goodsUUID"`
		OwnerUser string  `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID  string  `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		SerialNo  string  `gorm:"column:serial_no;comment:批次序号" json:"serialNo"`
		Price     float32 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
		Weight    float32 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
		Mount     int32   `gorm:"column:mount;comment:数量" json:"mount"`
		CustomerFeild
		GoodsFeild
	}
	// 订单操作记录
	BatchOrderOpr struct {
		BaseModel
		BatchOrderUID string `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		OwnerUser     string `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		UserUUID      string `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		History       Array  `gorm:"column:history;comment:历史" json:"History"`
	}
	BatchOrderPay struct {
		BaseModel
		CustomerUUID   string  `gorm:"column:customer_uuid;comment:" json:"customerUUID"`
		OwnerUser      string  `gorm:"column:owner_user;comment:所属用户" json:"ownerUser"`
		BatchOrderUUID string  `gorm:"column:batch_order_uuid;comment:批次uuid'" json:"batchOrderUUID"`
		Amount         float32 `gorm:"column:amount;type:decimal(10,2);comment:支付金额" json:"amount"`
		PayType        int32   `gorm:"column:pay_type;comment:付款方式" json:"payType"`
	}
)

func (bo *BatchOrder) UpdateStatus(db *gorm.DB, status int32) error {
	return WhereUIDCond(bo.UID).Cond(db).Model(&BatchOrder{}).Update("status", status).Error
}

func (bo *BatchOrder) UpdateShare(db *gorm.DB) error {
	return WhereUIDCond(bo.UID).Cond(db).Updates(&BatchOrder{
		Shared: BatchOrderShared,
	}).Error
}

func (bo *BatchOrder) Update(db *gorm.DB) error {
	return WhereUIDCond(bo.UID).Cond(db).Updates(&BatchOrder{
		Shared: BatchOrderShared,
	}).Error
}

func (b *BatchOrderGoods) Update(db *gorm.DB) error {
	return WhereUIDCond(b.UID).Cond(db).Updates(&BatchOrderGoods{
		Price:  b.Price,
		Mount:  b.Mount,
		Weight: b.Weight,
	}).Error
}

func (b *BatchOrderPay) Update(db *gorm.DB) error {
	return WhereUIDCond(b.UID).Cond(db).Updates(&BatchOrderPay{
		PayType: b.PayType,
		Amount:  b.Amount,
	}).Error
}

func (b *BatchOrderGoods) Amount() float32 {
	if b.Weight == 0 {
		return b.Price * float32(b.Mount)
	}
	return b.Price * b.Weight
}

func (b *BatchOrder) SetTotalAmount() float32 {
	var t float32
	for _, batchGoods := range b.GoodsListRelated {
		t += batchGoods.Amount()
	}
	return float32(math.Round(float64(t)))
}

func (b *BatchOrder) SetCreditAmount(pay float32) {
	b.CreditAmount -= pay
}
