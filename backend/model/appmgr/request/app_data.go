package request

import (
	"github.com/msean/csj/backend/model/common/request"
	"time"
)

type AppDataSearch struct {
	request.PageInfo
	OwnerUser    string     `form:"ownerUser" binding:"max=64"`
	CustomerUUID string     `form:"customerUUID" binding:"max=64"`
	BatchUUID    string     `form:"batchUUID" binding:"max=64"`
	Keyword      string     `form:"keyword" binding:"max=200"`
	Status       *int       `form:"status"`
	From         *time.Time `form:"from" time_format:"2006-01-02T15:04:05Z07:00"`
	To           *time.Time `form:"to" time_format:"2006-01-02T15:04:05Z07:00"`
}
