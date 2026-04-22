// 自动生成模板BanRecord
package bot

import (
	"time"

	"github.com/msean/csj/backend/global"
)

// 封禁记录 结构体  BanRecord
type BanRecord struct {
	global.GVA_MODEL
	BotID       int64      `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                      //机器人ID
	BotName     string     `json:"botName" form:"botName" gorm:"-"`                                             //机器人ID
	UserID      int64      `json:"userID" form:"userID" gorm:"comment:用户ID;column:user_id;"`                    //用户ID
	UserName    string     `json:"userName" form:"userName" gorm:"comment:用户名;column:user_name;size:128;index"` //userName
	FullName    string     `json:"fullName" form:"fullName" gorm:"comment:用户名;column:full_name;size:128;"`      //userName
	ChatID      int64      `json:"chatID" form:"chatID" gorm:"comment:所在群聊;column:chat_id;"`                    //chatID
	ChatName    string     `json:"chatName" form:"chatName" gorm:"-"`                                           //chatName
	BanDuration int64      `json:"banDuration" form:"banDuration" gorm:"comment:封禁时长;column:ban_duration;"`     //封禁时长
	BanType     int        `json:"banType" form:"banType" gorm:"comment:封禁时长;column:ban_type;"`                 //封禁类型 1 消息 2 成员
	Remark      string     `json:"reMark" form:"reMark" gorm:"comment:封禁时长;column:remark;"`                     //封禁备注
	Msg         string     `json:"msg" form:"msg" gorm:"type:text;comment:禁用信息;column:msg;"`                    //封禁消息
	LiftingTime *time.Time `json:"liftingTime" form:"liftingTime" gorm:"comment:解禁时间;column:lifting_time;"`     //解禁时间
	Status      int        `json:"status" form:"status" gorm:"-"`                                               //状态 1 封禁中 2 解禁
}

// TableName 封禁记录 BanRecord自定义表名 ban_record
func (BanRecord) TableName() string {
	return "ban_record"
}
