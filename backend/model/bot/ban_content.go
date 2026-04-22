// 自动生成模板BotMsgMgr
package bot

import (
	"github.com/msean/csj/backend/global"
)

// 机器人消息管理 结构体  BotBanContent
type BotBanContent struct {
	global.GVA_MODEL
	BanContent string `json:"banContent" form:"banContent" gorm:"comment:ban_content;column:ban_content;size:1024;"` //禁用内容
	BotID      int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                                //机器人ID
	BotName    string `json:"botName" form:"botName" gorm:"-"`                                                       //机器人ID
}

func (BotBanContent) TableName() string {
	return "bot_ban_content"
}
