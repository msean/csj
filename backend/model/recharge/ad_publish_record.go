// 自动生成模板AdPublishRecord
package recharge

import (
	"github.com/msean/csj/backend/global"
)

// 广告发布记录 结构体  AdPublishRecord
type AdPublishRecord struct {
	global.GVA_MODEL
	BotID        int64   `json:"botID" form:"botID" gorm:"column:bot_id;"`                             //机器人ID
	BotName      string  `json:"botName" form:"botName"`                                               //机器人ID
	ChannelName  string  `json:"channelName" form:"channelName"`                                       //机器人ID
	PublishTimes int     `json:"publishTimes" form:"publishTimes" gorm:"column:publish_times;"`        //发布次数
	UserID       int64   `json:"userID" form:"userID" gorm:"column:user_id;"`                          //发布用户ID
	UserName     string  `json:"userName" form:"userName" gorm:"column:user_name;size:256"`            //发布用户ID
	ChannelID    int64   `json:"channelID" form:"channelID" gorm:"column:channel_id;size:256"`         //发布用户ID
	Price        float64 `json:"price" form:"price" gorm:"column:price;type:decimal(10,3)"`            //发布价格
	Content      string  `json:"content" form:"content" gorm:"comment:发布内容;column:content;type:text;"` //发布内容
}

// TableName 广告发布记录 AdPublishRecord自定义表名 bot_ad_publish_record
func (AdPublishRecord) TableName() string {
	return "bot_ad_publish_record"
}
