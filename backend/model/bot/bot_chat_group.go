// 自动生成模板BotChatGroup
package bot

import (
	"time"

	"github.com/msean/csj/backend/global"
)

// 机器人群组列表 结构体  BotChatGroup
type BotChatGroup struct {
	global.GVA_MODEL
	BotName               string `json:"botName" form:"botName" gorm:"-"`
	BotID                 int64  `json:"botID" form:"botID" gorm:"column:bot_id;index:idx_bot_group,priority:1"`
	ChatGroupID           int64  `json:"chatGroupID" form:"chatGroupID" gorm:"column:chat_group_id;index:idx_bot_group,priority:2"`
	ChatGroupName         string `json:"chatGroupName" form:"chatGroupName" gorm:"column:chat_group_name;"`
	BanForward            int64  `json:"banForward" form:"banForward" gorm:"column:ban_forward;default:1"`
	MaxWords              int64  `json:"maxWords" form:"maxWords" gorm:"column:max_words;default:-1"`
	SyncMessage           int64  `json:"syncMessage" form:"botID" gorm:"column:sync_message;default:2"`
	MustJoinChannels      string `json:"mustJoinChannels" form:"mustJoinChannels" gorm:"column:must_join_channels"`
	InvaidChannelFoldLink string `json:"invaidChannelFoldLink" form:"invaidChannelFoldLink" gorm:"column:invaid_channel_fold_link"`
	Members               string `json:"members" form:"members" gorm:"column:members;type:text;"`
}

type BotChatGroupClassify struct {
	global.GVA_MODEL
	Title      string `json:"title" form:"title" gorm:"column:title"`                                         //标题
	ChatGroups string `json:"chatGroups" form:"chatGroups" gorm:"comment:群组列表;column:chat_groups;type:text;"` //机器人群组
	Users      string `json:"permitUsers" form:"permitUsers" gorm:"comment:群组列表;column:users;type:text;"`     //机器人群组
	Refresh    bool   `json:"refresh" form:"refresh" gorm:"-"`                                                // 允许操作用户
	// ChatGroupMapper map[int]string `json:"chatGroupMapper" form:"chatGroupMapper" gorm:"-"`
	// UserMapper      map[int]string `json:"userMapper" form:"userMapper" gorm:"-"`
}

type BotChatGroupRelatedChannelFollow struct {
	UserID      int64     `json:"userID" form:"userID" gorm:"column:user_id;index"`
	BotID       int64     `json:"botID" form:"botID" gorm:"column:bot_id"`
	ChatGroupID int64     `json:"chatGroupID" form:"chatGroupID" gorm:"column:chat_group_id"`
	CheckTime   time.Time `json:"checkTime" form:"checkTime" gorm:"column:check_time"`
}

// TableName 机器人群组列表 BotChatGroup自定义表名 bot_chat_group
func (BotChatGroupRelatedChannelFollow) TableName() string {
	return "bot_group_related_channnel_follow"
}

// TableName 机器人群组列表 BotChatGroup自定义表名 bot_chat_group
func (BotChatGroup) TableName() string {
	return "bot_chat_group"
}

// TableName 机器人群组列表 BotChatGroup自定义表名 bot_chat_group
func (BotChatGroupClassify) TableName() string {
	return "bot_chat_group_classify"
}
