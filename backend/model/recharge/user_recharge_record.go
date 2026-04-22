// 自动生成模板UserRechargeRecord
package recharge

import (
	"time"

	"github.com/msean/csj/backend/global"
)

// 用户充值记录 结构体  UserRechargeRecord
type UserRechargeRecord struct {
	global.GVA_MODEL
	BotID           int64     `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                                     //机器人ID
	PublishTimes    int       `json:"publishTimes" form:"publishTimes" gorm:"comment:发布次数;column:publish_times;default:1"`        //发布次数
	StartTime       time.Time `json:"startTime" form:"startTime" gorm:"comment:发布开始时间;column:start_time;"`                        //发布开始时间
	PublishInterval int64     `json:"publishInterval" form:"publishInterval" gorm:"comment:发布间隔;column:publish_interval;"`        //发布间隔
	PublishContent  string    `json:"publishContent" form:"publishContent" gorm:"comment:发布内容;column:publish_content;type:text;"` //发布内容
	Status          int64     `json:"status" form:"status" gorm:"comment:状态(1、创建 2、支付成功 3、支付超时失败);column:status;index"`           //状态
	MsgID           int       `json:"msgID" form:"msgID" gorm:"column:msg_id;"`                                                   //状态
	ChannelID       int64     `json:"channelID" form:"channelID" gorm:"column:channel_id"`                                        //状态
	BotName         string    `json:"botName" form:"botName" gorm:"-;"`
	UserName        string    `json:"userName" form:"userName" gorm:"column:user_name;index"`
	UserID          int64     `json:"userID" form:"userID" gorm:"column:user_id"`
	Price           float64   `json:"price" form:"price" gorm:"type:decimal(10,3);column:price"`
	ChatID          int64     `json:"chatID" form:"chatID" gorm:"column:chat_id"`
	TxID            string    `json:"txID" form:"txID" gorm:"size:256;column:tx_id"`                       // 交易ID
	PaymentAddr     string    `json:"paymentAddr" form:"paymentAddr" gorm:"size:256; column:payment_addr"` // 收款地址
}

// TableName 用户充值记录 UserRechargeRecord自定义表名 bot_user_recharge_record
func (UserRechargeRecord) TableName() string {
	return "bot_user_recharge_record"
}
