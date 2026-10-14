package request

import "github.com/msean/csj/backend/model/common/request"

// AppVersionSearch App版本搜索条件
type AppVersionSearch struct {
	Version string `json:"version" form:"version"`
	request.PageInfo
}
