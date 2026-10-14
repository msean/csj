package appmgr

import "github.com/msean/csj/backend/service"

type ApiGroup struct {
	AppVersionApi
}

var (
	appVersionService = service.ServiceGroupApp.AppMgrServiceGroup.AppVersionService
)
