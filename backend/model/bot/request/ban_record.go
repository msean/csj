package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type BanRecordSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	UserName       string      `json:"userName" form:"userName"`
	request.PageInfo
}
