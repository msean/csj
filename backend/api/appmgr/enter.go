package appmgr

import "github.com/msean/csj/backend/service"

type ApiGroup struct {
	AppVersionApi
	AppDataApi
	OfficialMessageApi
	MemberConfigApi
	RechargeOrderApi
	AppUserApi
}

var (
	appDataService         = service.ServiceGroupApp.AppMgrServiceGroup.AppDataService
	officialMessageService = service.ServiceGroupApp.AppMgrServiceGroup.OfficialMessageService
	appVersionService      = service.ServiceGroupApp.AppMgrServiceGroup.AppVersionService
	memberConfigService    = service.ServiceGroupApp.AppMgrServiceGroup.MemberConfigService
	rechargeOrderService   = service.ServiceGroupApp.AppMgrServiceGroup.RechargeOrderService
	appUserService         = service.ServiceGroupApp.AppMgrServiceGroup.AppUserService
)
