package appmgr

import "github.com/msean/csj/backend/service"

type ApiGroup struct {
	AppVersionApi
	MemberConfigApi
	RechargeOrderApi
	AppUserApi
}

var (
	appVersionService    = service.ServiceGroupApp.AppMgrServiceGroup.AppVersionService
	memberConfigService  = service.ServiceGroupApp.AppMgrServiceGroup.MemberConfigService
	rechargeOrderService = service.ServiceGroupApp.AppMgrServiceGroup.RechargeOrderService
	appUserService       = service.ServiceGroupApp.AppMgrServiceGroup.AppUserService
)
