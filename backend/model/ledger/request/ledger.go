package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type LedgerSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}
