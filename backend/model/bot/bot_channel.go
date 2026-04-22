// 自动生成模板BotChannel
package bot

import (
	"github.com/msean/csj/backend/global"
)

// 机器人渠道 结构体  BotChannel
type BotChannel struct {
	global.GVA_MODEL
	BotID       int64  `json:"botID" form:"botID" gorm:"comment:机器人id;column:bot_id;"`                           //机器人id
	ChannelID   int64  `json:"channelID" form:"channelID" gorm:"comment:频道ID;column:channel_id;"`                //频道ID
	ChannelName string `json:"channelName" form:"channelName" gorm:"comment:渠道名称;column:channel_name;size:256;"` //渠道名称
	BotName     string `json:"botName" form:"botName" gorm:"-"`                                                  //渠道名称
}

// TableName 机器人渠道 BotChannel自定义表名 bot_channel
func (BotChannel) TableName() string {
	return "bot_channel"
}
