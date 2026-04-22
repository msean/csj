package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotBanContentSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	BanContent     string      `json:"banContent" form:"banContent"`
	BotID          int         `json:"botID" form:"botID"`
	request.PageInfo
}
