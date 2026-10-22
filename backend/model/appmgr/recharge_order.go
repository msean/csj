package appmgr

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RechargeOrder 会员充值订单
// 对应 app 客户端项目（app/service/model/recharge_order.go）的 recharge_orders 表，主键为 uid
type RechargeOrder struct {
	UID           string         `json:"uid" gorm:"column:uid;type:varchar(64);primaryKey"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index"`
	UserUUID      string         `json:"userUUID" gorm:"column:user_uuid;type:varchar(64);index;comment:用户UUID"`
	OrderNo       string         `json:"orderNo" gorm:"column:order_no;type:varchar(64);uniqueIndex;comment:商户订单号"`
	MemberType    int            `json:"memberType" gorm:"column:member_type;comment:会员类型 1月 2季度 3年"`
	Amount        float64        `json:"amount" gorm:"column:amount;type:decimal(10,2);comment:充值金额"`
	PayStatus     int            `json:"payStatus" gorm:"column:pay_status;default:0;index;comment:支付状态 0待支付 1已支付 2已取消 3已退款"`
	PayTime       *time.Time     `json:"payTime" gorm:"column:pay_time;comment:支付时间"`
	TransactionID string         `json:"transactionID" gorm:"column:transaction_id;type:varchar(64);comment:微信支付订单号"`
}

func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

func (m *RechargeOrder) BeforeCreate(tx *gorm.DB) error {
	if m.UID == "" {
		m.UID = uuid.NewString()
	}
	return nil
}
