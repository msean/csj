package model

import (
	"app/pkg/utils"
	"app/service/common"
	"time"
)

type (
	// 一个订单对应多个订单货品
	BatchOrder struct {
		BaseModel
		BatchUUID        string             `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		OwnerUser        string             `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID         string             `gorm:"column:user_uuid;comment:开单uuid;index:idx_owner_user" json:"customerUUID"`
		Shared           int                `gorm:"column:shared;comment:是否分享单" json:"shared"`
		SharedTime       int                `gorm:"column:share_time;comment:分享时间" json:"sharedTime"`
		Status           int                `gorm:"column:status;comment:状态" json:"status"`
		TotalAmount      float64            `gorm:"column:amount;comment:金额" json:"totalAmount"`
		CreditAmount     float64            `gorm:"column:credit_amount;comment:赊欠金额" json:"creditAmount"`
		GoodsListRelated []*BatchOrderGoods `gorm:"foreignKey:BatchOrderUID;references:UID" json:"goodsList"`
		CustomerFeild
	}
	BatchOrderGoods struct {
		BaseModel
		BatchOrderUID string  `gorm:"column:batch_order_uuid;comment:批次uuid" json:"batchOrderUUID"`
		BatchUUID     string  `gorm:"column:batch_uuid;comment:批次uuid" json:"batchUUID"`
		GoodsUUID     string  `gorm:"column:goods_uuid;comment:货品uuid" json:"goodsUUID"`
		OwnerUser     string  `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		UserUUID      string  `gorm:"column:user_uuid;comment:开单uuid" json:"customerUUID"`
		GoodType      int     `gorm:"column:good_type;comment:批次序号" json:"-"`
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
		CustomerUUID   string     `gorm:"column:customer_uuid;comment:客户UUID" json:"customerUUID"`
		OwnerUser      string     `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		BatchOrderUUID string     `gorm:"column:batch_order_uuid;comment:订单UUID（快捷还款为空）" json:"batchOrderUUID"`
		Amount         float64    `gorm:"column:amount;type:decimal(10,2);comment:还款金额（正数还款，负数撤销）" json:"amount"`
		PayType        int32      `gorm:"column:pay_type;comment:付款方式 1-现金 2-微信 3-支付宝" json:"payType"`
		Remark         string     `gorm:"column:remark;comment:备注" json:"remark"`
		IsRevoked      int        `gorm:"column:is_revoked;default:0;comment:是否已撤销 0-否 1-是" json:"isRevoked"`
		RevokedAt      *time.Time `gorm:"column:revoked_at;comment:撤销时间" json:"revokedAt"`
		RevokedReason  string     `gorm:"column:revoked_reason;comment:撤销原因" json:"revokedReason"`
		// 快捷还款时存储多个订单的还款详情 JSON
		// 格式：[{"orderUUID":"xxx","amount":100},{"orderUUID":"yyy","amount":200}]
		PayDetails string `gorm:"column:pay_details;type:text;comment:还款详情JSON" json:"payDetails"`
	}

	// MessageCenter 消息中心表
	MessageCenter struct {
		BaseModel
		OwnerUser    string `gorm:"column:owner_user;comment:所属用户;index" json:"ownerUser"`
		CustomerUUID string `gorm:"column:customer_uuid;comment:客户UUID" json:"customerUUID"`
		Type         int    `gorm:"column:type;comment:消息类型 1-还款 2-撤销还款" json:"type"`
		Event        string `gorm:"column:event;comment:事件名称" json:"event"`
		Content      string `gorm:"column:content;type:text;comment:消息内容" json:"content"`
		IsRead       int    `gorm:"column:is_read;default:0;comment:是否已读 0-未读 1-已读" json:"isRead"`
		RelatedUUID  string `gorm:"column:related_uuid;comment:关联UUID（还款记录UUID）" json:"relatedUUID"`
	}

	// PayDetail 还款详情（用于快捷还款记录多个订单）
	PayDetail struct {
		OrderUUID string  `json:"orderUUID"`
		Amount    float64 `json:"amount"`
	}
)

func (b *BatchOrder) DefaultSet() {
	b.Status = common.BatchOrderTemp
	b.Shared = common.BatchOrderUnshare
}

func (b *BatchOrderGoods) Amount() float64 {
	// 件数量为0 则是散装
	if b.GoodType == common.GoodsTypeBulk {
		return b.Price * float64(b.Weight)
	}
	return b.Price * float64(b.Mount)
}

func (b *BatchOrderGoods) Sell() float64 {
	// 件数量为0 则是散装
	if b.GoodType == common.GoodsTypeBulk {
		return b.Weight
	}
	return float64(b.Mount)
}

// func (b *BatchOrderGoods) GType() int {
// 	// 件数量为0 则是散装
// 	if b.Mount == 0 {
// 		return common.GoodsTypeBulk
// 	}
// 	return common.GoodsTypeFix
// }

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
	if b.GoodType == common.GoodsTypeBulk {
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
