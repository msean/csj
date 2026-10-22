package service

import (
	"app/global"
	"app/pkg/wechat"
	"app/service/cache"
	"app/service/handler"
	"app/service/model"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(cfgFilepath string, enableMigrate bool, enableShardMigrate bool) {
	err := global.InitRunTime(cfgFilepath)
	if err != nil {
		panic(err)
	}
	handler.InitEngine(global.Global.Engine)

	// 微信支付初始化（配置不完整时仅警告，不阻断启动）
	if err = wechat.Init(); err != nil {
		global.Global.Logger.Error("微信支付初始化失败", zap.Error(err))
	}

	// 根据命令行参数决定是否执行迁移
	if enableMigrate {
		model.Migrate()
		global.Global.Logger.Info("数据库迁移已执行")
	} else {
		global.Global.Logger.Info("数据库迁移已跳过（使用 -m 或 --migrate 参数启用）")
	}

	// 初始化分表：创建分表并迁移现有数据
	if enableShardMigrate {
		if err := model.MigrateShardingTables(); err != nil {
			global.Global.Logger.Error("创建分表失败", zap.Error(err))
			panic(err)
		}

		// 迁移现有数据到分表（仅首次执行）
		// if err := model.MigrateShardingData(); err != nil {
		// 	global.Global.Logger.Error("迁移数据到分表失败", zap.Error(err))
		// 	panic(err)
		// }
		global.Global.Logger.Info("分表迁移已执行")
	} else {
		global.Global.Logger.Info("分表迁移已跳过（使用 -s 或 --shard-migrate 参数启用）")
	}

	cache.Init()

	defer global.Global.Close()
	go func() { global.Global.Run(global.Global.Engine) }()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Global.Close()
}
