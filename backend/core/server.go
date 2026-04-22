package core

import (
	"fmt"
	"time"

	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/initialize"
	"github.com/msean/csj/backend/service/system"
)

func RunServer() {
	if global.GVA_CONFIG.System.UseRedis {
		initialize.Redis()
		if global.GVA_CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}

	if global.GVA_MYSQL != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)

	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
