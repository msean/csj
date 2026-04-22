package ledger2

import (
	"github.com/msean/csj/backend/global"
)

// 帐薄 结构体  Ledger
type Ledger struct {
	global.GVA_MODEL
	OprUserID        int64   `json:"oprUserID" form:"oprUserID" gorm:"comment:操作用户ID;column:opr_user_id;"`                   //操作用户ID
	OprUserFirstName string  `json:"oprUsername" form:"oprUsername" gorm:"column:opr_first_name;"`                           //操作人的用户名称
	OprUserLastName  string  `json:"OprUserLastName" form:"OprUserLastName" gorm:"column:opr_last_name;"`                    //操作人的用户名称
	OprUserNickname  string  `json:"oprUserNickname" form:"oprUserNickname" gorm:"comment:操作人昵称;column:opr_user_nick_name;"` //操作人昵称
	ActionType       int     `json:"actionType" form:"actionType" gorm:"comment:操作类型;column:action_type;"`                   //1 收入 2 支出
	Amount           float64 `json:"amount" form:"amount" gorm:"comment:操作金额;type:decimal(18,2);column:amount;"`             //金额
	UserName         string  `json:"userName" form:"userName" gorm:"comment:姓名;column:user_name;"`                           // botID
	Account          string  `json:"account" form:"account" gorm:"comment:账号;column:account;"`                               // 姓名
	MessageID        int     `gorm:"uniqueIndex:uniq_msg;column:message_id;comment:消息ID"`
	RawInput         string  `json:"rawInput" form:"rawInput" gorm:"comment:原始输入;column:raw_input;"` //原始输入
	Remark           string  `json:"remark" form:"remark" gorm:"comment:备注;column:remark;"`
	BotID            int64   `gorm:"index:idx_bot_chat_date,priority:1;column:bot_id;comment:机器人ID"`
	ChatGroupID      int64   `gorm:"index:idx_bot_chat_date,priority:2;column:chat_group_id;comment:群ID"`
	WorkDate         string  `gorm:"index:idx_bot_chat_date,priority:3;column:work_date;comment:账单日期"`
}

type LedgerPermission struct {
	global.GVA_MODEL
	// CurrentFeeRate float64 `json:"currentFeeRate" form:"currentFeeRate" gorm:"comment:当前费率;column:current_fee_rate"`
	BotName       string `json:"botName" form:"botName" gorm:"-;"`
	ChatGroupName string `json:"chatGroupName" form:"chatGroupName" gorm:"-;"`
	BotID         int64  `json:"botID" form:"botID" gorm:"comment:机器人ID;column:bot_id"`
	ChatGroupID   int64  `json:"chatGroupID" form:"chatGroupID" gorm:"comment:群聊ID;column:chat_group_id;"` //群聊ID
	OprUsers      string `json:"oprUsers" form:"oprUsers" gorm:"comment:操作人;column:opr_users;type:text;"`  //操作人
}

type LedgerSession struct {
	global.GVA_MODEL
	ID          uint   `gorm:"primaryKey"`
	BotID       int64  `gorm:"index;column:bot_id"`
	ChatGroupID int64  `gorm:"index;column:chat_group_id"`
	WorkDate    string `gorm:"index;column:work_date"` // 2026-03-24
	IsActive    int    `gorm:"column:is_active;default:0"`
}

// TableName 帐薄 Ledger自定义表名 ledger
func (Ledger) TableName() string {
	return "bot_ledger2"
}

// TableName 帐薄权限管理 LedgerPermission自定义表名 ledger_permission
func (LedgerPermission) TableName() string {
	return "bot_ledger2_permission"
}
