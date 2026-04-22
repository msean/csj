package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotBanGroupMemSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	BotID          int         `json:"botID" form:"botID"`
	request.PageInfo
}
