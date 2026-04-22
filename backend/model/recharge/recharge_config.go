// 自动生成模板RechargeConfig
package recharge

import (
	"github.com/msean/csj/backend/global"
)

// 充值配置 结构体  RechargeConfig
type RechargeConfig struct {
	global.GVA_MODEL
	BotID        int64   `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`                              //机器人ID
	PublishTimes int64   `json:"publishTimes" form:"publishTimes" gorm:"comment:发布次数;column:publish_times;default:1"` //发布次数
	Price        float64 `json:"price" form:"price" gorm:"comment:价格;column:price;"`                                  //价格
	BotName      string  `json:"botName" form:"botName" gorm:"-;"`
}

// TableName 充值配置 RechargeConfig自定义表名 bot_recharge_config
func (RechargeConfig) TableName() string {
	return "bot_recharge_config"
}
