package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/csj_customers"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/user"
)

type ServiceGroup struct {
	SystemServiceGroup        system.ServiceGroup
	ExampleServiceGroup       example.ServiceGroup
	UserServiceGroup          user.ServiceGroup
	Csj_customersServiceGroup csj_customers.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
