// 自动生成模板Ledger
package ledger

import (
	"github.com/msean/csj/backend/global"
)

// 帐薄 结构体  Ledger
type Ledger struct {
	global.GVA_MODEL
	BotName          string  `json:"botName" form:"botName" gorm:"-;"`
	ChatGroupName    string  `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`
	OprUserID        int64   `json:"oprUserID" form:"oprUserID" gorm:"comment:操作用户ID;column:opr_user_id;"`                                  //操作用户ID
	OprUserFirstName string  `json:"oprUsername" form:"oprUsername" gorm:"column:opr_first_name;"`                                          //操作人的用户名称
	OprUserLastName  string  `json:"OprUserLastName" form:"OprUserLastName" gorm:"column:opr_last_name;"`                                   //操作人的用户名称
	OprUserNickname  string  `json:"oprUserNickname" form:"oprUserNickname" gorm:"comment:操作人昵称;column:opr_user_nick_name;"`                //操作人昵称
	ActionType       int     `json:"actionType" form:"actionType" gorm:"comment:操作类型;column:action_type;"`                                  //操作类型 1 入款 2 下发
	Amount           float64 `json:"amount" form:"amount" gorm:"comment:操作金额;type:decimal(32,2);column:amount;"`                            //操作金额
	AmountWithFee    float64 `json:"amount_with_fee" form:"amount_with_fee" gorm:"comment:操作金额;type:decimal(32,2);column:amount_with_fee;"` //操作金额
	BotID            int64   `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id;"`
	ChatGroupID      int64   `json:"chatGroupID" form:"chatGroupID" gorm:"comment:所在群组;column:chat_group_id;"` //所在群组
	MessageID        int64   `json:"messageID" form:"messageID" gorm:"index;comment:消息ID;column:message_id;"`  //消息ID
	RawInput         string  `json:"rawInput" form:"rawInput" gorm:"comment:原始输入;column:raw_input;"`           //原始输入
	Remark           string  `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`                    //原始输入
	CurrentFeeRate   float64 `json:"currentFeeRate" form:"currentFeeRate" gorm:"comment:当前费率;column:current_fee_rate"`
}

type LedgerPermission struct {
	global.GVA_MODEL
	CurrentFeeRate float64 `json:"currentFeeRate" form:"currentFeeRate" gorm:"comment:当前费率;column:current_fee_rate"`
	BotName        string  `json:"botName" form:"botName" gorm:"-;"`
	ChatGroupName  string  `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`
	BotID          int64   `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id"`
	ChatGroupID    int64   `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊ID
	OprUsers       string  `json:"oprUsers" form:"oprUsers" gorm:"comment:操作人;column:opr_users;type:text;"`  //操作人
}

// TableName 帐薄 Ledger自定义表名 ledger
func (Ledger) TableName() string {
	return "bot_ledger"
}

// TableName 帐薄权限管理 LedgerPermission自定义表名 ledger_permission
func (LedgerPermission) TableName() string {
	return "bot_ledger_permission"
}
