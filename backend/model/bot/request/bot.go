package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BotSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Name           string      `json:"name" form:"name"`
	Token          string      `json:"token" form:"token"`
	request.PageInfo
}
