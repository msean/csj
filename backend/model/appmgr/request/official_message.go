package request

import "github.com/msean/csj/backend/model/common/request"

type OfficialMessageSearch struct {
	Status *int `json:"status" form:"status" binding:"omitempty,oneof=0 1"`
	request.PageInfo
}
