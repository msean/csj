package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type UserRechargeRecordSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	BotID          int64       `json:"botID" form:"botID"`
	Status         int64       `json:"status" form:"status"`
	UserID         int64       `json:"userID" form:"userID"`
	request.PageInfo
}
