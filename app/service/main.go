package service

import (
	"app/global"
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
	handler.InitEngine(global.GlobalRunTime.Engine)
	model.Migrate()

	defer global.GlobalRunTime.Close()
	go func() { global.GlobalRunTime.Run(global.GlobalRunTime.Engine) }()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.GlobalRunTime.Close()
}
