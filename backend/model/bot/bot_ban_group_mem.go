// 自动生成模板BotBanGroupMem
package bot

import (
	"github.com/msean/csj/backend/global"
)

// 封禁成员设置 结构体  BotBanGroupMem
type BotBanGroupMem struct {
	global.GVA_MODEL
	BotID         int64  `json:"botID" form:"botID" gorm:"column:bot_id;"` //botID
	BotName       string `json:"botName" form:"botName" gorm:"-;"`
	ChatGroupName string `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`                                                          //botID
	ChatGroupID   int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:chatGroupID;column:chat_group_id;"`                       //chatGroupID
	BanMemContent string `json:"banMemContent" form:"banMemContent" gorm:"comment:封禁成员名称(成员还有该字段就封禁);column:ban_mem_content;size:256;"` //banMemContent
}

// TableName 封禁成员设置 BotBanGroupMem自定义表名 bot_ban_group_mem
func (BotBanGroupMem) TableName() string {
	return "bot_ban_group_mem"
}
