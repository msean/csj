package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/csj_customers"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/user"
)

type RouterGroup struct {
	System        system.RouterGroup
	Example       example.RouterGroup
	User          user.RouterGroup
	Csj_customers csj_customers.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
