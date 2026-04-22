package public

import (
	api "github.com/msean/csj/backend/api"
)

type RouterGroup struct {
	PublicRouter
}

var (
	medioApi = api.ApiGroupApp.PublicApiGroup.MedioApi
)
