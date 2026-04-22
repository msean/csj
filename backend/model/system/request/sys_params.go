package request

import (
	"time"

	"github.com/msean/csj/backend/model/common/request"
)

type SysParamsSearch struct {
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Name           string     `json:"name" form:"name" `
	Key            string     `json:"key" form:"key" `
	request.PageInfo
}
