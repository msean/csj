package request

import "github.com/msean/csj/backend/model/common/request"

// AppUserSearch 用户搜索条件
type AppUserSearch struct {
	Phone string `json:"phone" form:"phone"`
	Name  string `json:"name" form:"name"`
	request.PageInfo
}
