// 自动生成模板UserWallet
package recharge

import (
	"github.com/msean/csj/backend/global"
)

// 用户钱包 结构体  UserWallet
type UserWallet struct {
	global.GVA_MODEL
	BotID    int64   `json:"botID" form:"botID" gorm:"comment:机器人名称;column:bot_id;uniqueIndex:uk_bot_user"`    //机器人ID
	BotName  string  `json:"botName" form:"botName"`                                                           //机器人ID
	UserID   int64   `json:"userID" form:"userID" gorm:"comment:用户ID;column:user_id;uniqueIndex:uk_bot_user;"` //用户ID
	UserName string  `json:"userName" form:"userName" gorm:"comment:用户名;column:user_name;"`                    //用户名称
	Balance  float64 `json:"balance" form:"balance" gorm:"column:balance;type:decimal(10,3);comment:余额"`       //余额
}

// TableName 用户钱包 UserWallet自定义表名 bot_user_wallet
func (UserWallet) TableName() string {
	return "bot_user_wallet"
}
