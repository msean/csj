package appmgr

import api "github.com/msean/csj/backend/api"

type RouterGroup struct {
	AppVersionRouter
	MemberConfigRouter
	RechargeOrderRouter
	AppUserRouter
}

var (
	appVersionApi    = api.ApiGroupApp.AppMgrApiGroup.AppVersionApi
	memberConfigApi  = api.ApiGroupApp.AppMgrApiGroup.MemberConfigApi
	rechargeOrderApi = api.ApiGroupApp.AppMgrApiGroup.RechargeOrderApi
	appUserApi       = api.ApiGroupApp.AppMgrApiGroup.AppUserApi
)
