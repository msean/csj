package request

import "github.com/msean/csj/backend/model/common/request"

// MemberConfigSearch 会员收费配置搜索条件
type MemberConfigSearch struct {
	Typ int `json:"type" form:"type"`
	request.PageInfo
}
