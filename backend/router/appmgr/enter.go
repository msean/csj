package appmgr

import api "github.com/msean/csj/backend/api"

type RouterGroup struct {
	AppVersionRouter
	AppDataRouter
	OfficialMessageRouter
	MemberConfigRouter
	RechargeOrderRouter
	AppUserRouter
}

var (
	appDataApi         = api.ApiGroupApp.AppMgrApiGroup.AppDataApi
	officialMessageApi = api.ApiGroupApp.AppMgrApiGroup.OfficialMessageApi
	appVersionApi      = api.ApiGroupApp.AppMgrApiGroup.AppVersionApi
	memberConfigApi    = api.ApiGroupApp.AppMgrApiGroup.MemberConfigApi
	rechargeOrderApi   = api.ApiGroupApp.AppMgrApiGroup.RechargeOrderApi
	appUserApi         = api.ApiGroupApp.AppMgrApiGroup.AppUserApi
)
