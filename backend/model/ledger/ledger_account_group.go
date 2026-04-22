// 自动生成模板LedgerAccountGroup
package ledger

import (
	"github.com/msean/csj/backend/global"
)

// 记账账号组 结构体  LedgerAccountGroup
type LedgerAccountGroup struct {
	global.GVA_MODEL
	AccountGroup *string `json:"accountGroup" form:"accountGroup" gorm:"column:account_group;type:text;"` //账号分组
	Title        *string `json:"title" form:"title" gorm:"comment:标题;column:title;size:256;"`             //tilte
}

// TableName 记账账号组 LedgerAccountGroup自定义表名 ledger_account_group
func (LedgerAccountGroup) TableName() string {
	return "ledger_account_group"
}
