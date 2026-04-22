package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type LedgerPermissionSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	request.PageInfo
}
