// 自动生成模板BotCmdConfig
package bot

import (
	"encoding/json"

	"github.com/msean/csj/backend/global"
)

// 机器人命令配置 结构体  BotCmdConfig
type BotCmdConfig struct {
	global.GVA_MODEL
	BotID      int64           `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                           //机器人ID
	Title      string          `json:"title" form:"title" gorm:"comment:标题;column:title"`                                //开始设置内容
	BotName    string          `json:"botName" form:"botName" gorm:"-"`                                                  //机器人ID
	Content    string          `json:"content" form:"content" gorm:"comment:开始设置内容;column:content;type:text;"`           //开始设置内容
	Cmd        string          `json:"cmd" form:"cmd" gorm:"comment:开始设置内容;column:cmd;;size:512;default:'/start'"`       //开始设置内容
	CmdButtons json.RawMessage `json:"cmdButtons" form:"cmdButtons" gorm:"comment:命令按钮配置;column:cmd_buttons;type:text;"` //命令按钮配置
	Type       int             `json:"type" form:"type" gorm:"comment:类型;column:type;default:1"`                         //开始设置内容
}

// TableName 机器人命令配置 BotCmdConfig自定义表名 bot_cmd_config
func (BotCmdConfig) TableName() string {
	return "bot_cmd_config"
}
