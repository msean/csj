package service

import (
	"app/global"
	"app/service/cache"
	"app/service/handler"
	"app/service/model"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfgFilepath string) {
	err := global.InitRunTime(cfgFilepath)
	if err != nil {
		panic(err)
	}
	handler.InitEngine(global.Global.Engine)
	model.Migrate()
	cache.Init()

	defer global.Global.Close()
	go func() { global.Global.Run(global.Global.Engine) }()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Global.Close()
}
