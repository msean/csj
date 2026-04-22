package response

import (
	"github.com/msean/csj/backend/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
