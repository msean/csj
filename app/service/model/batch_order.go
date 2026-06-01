package model

import (
	"app/pkg/utils"
	"app/service/common"

	"gorm.io/gorm"
)

type (
	// 一个订单对应多个订单货品
	BatchOrder struct {
		BaseModel
		BatchUUID        string             `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		OwnerUser        string             `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID         string             `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		Shared           int                `gorm:"column:shared;comment:是否分享单" json:"shared"`
		SharedTime       int                `gorm:"column:share_time;comment:分享时间" json:"sharedTime"`
		Status           int                `gorm:"column:status;comment:状态" json:"status"`
		TotalAmount      float64            `gorm:"column:amount;comment:金额" json:"totalAmount"`
		CreditAmount     float64            `gorm:"column:credit_amount;comment:赊欠金额" json:"creditAmount"`
		GoodsListRelated []*BatchOrderGoods `gorm:"foreignKey:BatchOrderUID;references:UID" json:"goodsList"`
		CustomerFeild
	}
	// todo,假如是散装一定要将Mount置为0
	BatchOrderGoods struct {
		BaseModel
		BatchOrderUID string  `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		BatchUUID     string  `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		GoodsUUID     string  `gorm:"column:goods_uuid;comment:货品uuid" json:"goodsUUID"`
		OwnerUser     string  `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID      string  `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		SerialNo      string  `gorm:"column:serial_no;comment:批次序号" json:"serialNo"`
		Price         float64 `gorm:"column:price;type:decimal(10,2);comment:单价" json:"price"`
		Weight        float64 `gorm:"column:weight;type:decimal(10,2);comment:重量" json:"weight"`
		Mount         int     `gorm:"column:mount;comment:数量" json:"mount"`                    // 定装 件数 散装 斤数
		Total         float64 `gorm:"column:total;type:decimal(10,2);comment:小计" json:"total"` // 小记
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
		Amount         float64 `gorm:"column:amount;type:decimal(10,2);comment:支付金额" json:"amount"`
		PayType        int32   `gorm:"column:pay_type;comment:付款方式" json:"payType"`
	}
)

func (bo *BatchOrder) UpdateStatus(db *gorm.DB, status int) error {
	return WhereUIDCond(bo.UID).Cond(db).Model(&BatchOrder{}).Update("status", status).Error
}

func (bo *BatchOrder) UpdateShare(db *gorm.DB) error {
	return WhereUIDCond(bo.UID).Cond(db).Updates(&BatchOrder{
		Shared: common.BatchOrderShared,
	}).Error
}

func (bo *BatchOrder) Update(db *gorm.DB) error {
	return WhereUIDCond(bo.UID).Cond(db).Updates(&BatchOrder{
		Shared: common.BatchOrderShared,
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

func (b *BatchOrderGoods) Amount() float64 {
	// 件数量为0 则是散装
	if b.Mount == 0 {
		return b.Price * float64(b.Weight)
	}
	return b.Price * float64(b.Mount)
}

func (b *BatchOrderGoods) GType() int {
	// 件数量为0 则是散装
	if b.Mount == 0 {
		return common.GoodsTypeBulk
	}
	return common.GoodsTypeFix
}

// func (b *BatchOrder) SetTotalAmount() float64 {
// 	var t float64
// 	for _, batchGoods := range b.GoodsListRelated {
// 		t += batchGoods.Amount()
// 	}
// 	return float64(math.Round(float64(t)))
// }

// func (b *BatchOrder) SetCreditAmount(pay float64) {
// 	b.CreditAmount -= pay
// }

// 设置小计件
func (b *BatchOrderGoods) SetTotal() float64 {
	if b.Mount == 0 {
		b.Total = utils.FloatReserve(b.Price*b.Weight, 0)
	} else {
		b.Total = utils.FloatReserve(float64(b.Mount)*b.Price, 0)
	}
	return b.Total
}

// func (bg *BatchOrderGoods) Amount() (amount string) {
// 	if bg.Type == common.GoodsTypeFix {
// 		amount = fmt.Sprintf("%d", bg.Mount)
// 		return
// 	}
// 	amount = fmt.Sprintf("%.2f", bg.Mount)
// 	return
// }
