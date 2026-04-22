package response

import "github.com/msean/csj/backend/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
