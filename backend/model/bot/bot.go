// 自动生成模板Bot
package bot

import (
	"time"

	"gorm.io/gorm"
)

// 机器人 结构体  Bot
type Bot struct {
	CreatedAt    time.Time      `json:"createdAt" form:"createdAt" gorm:"column:created_at;"`
	UpdatedAt    time.Time      `json:"updatedAt" form:"updatedAt" gorm:"column:updated_at;"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `json:"name" form:"name" gorm:"comment:机器人名称;column:name;"`
	IsForLedger  int64          `json:"isForLedger" form:"isForLedger" gorm:"comment:是否用来记账(1开2关);column:is_for_ledger;default:2"`     // 是否用来做记账
	IsForLedger2 int64          `json:"isForLedger2" form:"isForLedger2" gorm:"comment:是否用来记账(1开2关);column:is_for_ledger2;default:2"`  // 是否用来做记账
	IsForMsgMgr  int64          `json:"isForMsgMgr" form:"isForMsgMgr" gorm:"comment:是否用来记账(1开2关);column:is_for_msg_mgr;default:2"`    // 是否用来做消息管理
	IsForMsgMass int64          `json:"isForMsgMass" form:"isForMsgMass" gorm:"comment:是否用来记账(1开2关);column:is_for_msg_mass;default:2"` // 是否用来做群发消息
	IsAdPublish  int64          `json:"isAdPublish" form:"isAdPublish" gorm:"comment:是否用来广告自动发布;column:is_for_ad_publish;default:2"`   // 是否用来广告自动发布
	BotID        int64          `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;index"`
	Token        string         `json:"token" form:"token" gorm:"comment:token;column:token;size:256;"`
	Chats        []BotChatGroup `json:"botChatGroups" form:"botChatGroups" gorm:"foreignKey:BotID;references:BotID"`
	Channels     []BotChannel   `json:"botChannels" form:"botChannels" gorm:"foreignKey:BotID;references:BotID"`
}

// TableName 机器人 Bot自定义表名 bot
func (Bot) TableName() string {
	return "bot"
}

type BotFeildExtend struct {
	BotName       string `json:"botName" form:"botName" gorm:"-"` //渠道名称
	ChatGroupName string `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`
}
