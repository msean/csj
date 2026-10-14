package service

import (
	"github.com/msean/csj/backend/service/appmgr"
	"github.com/msean/csj/backend/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
	AppMgrServiceGroup appmgr.ServiceGroup
}
