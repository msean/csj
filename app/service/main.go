package service

import (
	"app/global"
	"app/service/handler"
	"app/service/model"
)

func Run(cfgFilepath string) {
	err := global.InitRunTime(cfgFilepath)
	if err != nil {
		panic(err)
	}
	handler.InitEngine(global.GlobalRunTime.Engine)
	model.Migrate()
	defer global.GlobalRunTime.Close()

	global.GlobalRunTime.Run()
}
