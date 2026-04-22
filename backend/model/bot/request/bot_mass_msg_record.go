package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotMassMsgRecordSearch struct {
	BotID          *int64      `json:"botID" form:"botID"`                     // 可选：机器人ID
	ChatGroupID    *int64      `json:"chatGroupID" form:"chatGroupID"`         // 可选：群聊ID
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"` // 创建日期范围
	request.PageInfo
}
