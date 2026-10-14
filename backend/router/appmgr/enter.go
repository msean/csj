package appmgr

import api "github.com/msean/csj/backend/api"

type RouterGroup struct {
	AppVersionRouter
}

var (
	appVersionApi = api.ApiGroupApp.AppMgrApiGroup.AppVersionApi
)
