package initialize

import (
	"os"

	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/bot"
	"github.com/msean/csj/backend/model/ledger"
	"github.com/msean/csj/backend/model/ledger2"
	"github.com/msean/csj/backend/model/recharge"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm(dbType string) *gorm.DB {
	switch dbType {
	case "mysql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	case "pgsql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Pgsql.Dbname
		return GormPgSql()
	case "oracle":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Oracle.Dbname
		return GormOracle()
	case "mssql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mssql.Dbname
		return GormMssql()
	case "sqlite":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Sqlite.Dbname
		return GormSqlite()
	default:
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		return GormMysql()
	}
}

func RegisterTables() {
	if global.GVA_CONFIG.System.DisableAutoMigrate {
		global.GVA_LOG.Info("auto-migrate is disabled, skipping table registration")
		return
	}

	// db := global.GVA_MYSQL
	// err := db.AutoMigrate(

	// 	system.SysApi{},
	// 	system.SysIgnoreApi{},
	// 	system.SysUser{},
	// 	system.SysBaseMenu{},
	// 	system.JwtBlacklist{},
	// 	system.SysAuthority{},
	// 	system.SysDictionary{},
	// 	system.SysOperationRecord{},
	// 	system.SysAutoCodeHistory{},
	// 	system.SysDictionaryDetail{},
	// 	system.SysBaseMenuParameter{},
	// 	system.SysBaseMenuBtn{},
	// 	system.SysAuthorityBtn{},
	// 	system.SysAutoCodePackage{},
	// 	system.SysExportTemplate{},
	// 	system.Condition{},
	// 	system.JoinTemplate{},
	// 	system.SysParams{},
	// 	system.SysVersion{},

	// 	example.ExaFile{},
	// 	example.ExaCustomer{},
	// 	example.ExaFileChunk{},
	// 	example.ExaFileUploadAndDownload{},
	// 	example.ExaAttachmentCategory{},
	// )
	// if err != nil {
	// 	global.GVA_LOG.Error("register table failed", zap.Error(err))
	// 	os.Exit(0)
	// }

	err := bizModel()

	if err != nil {
		global.GVA_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}

func bizModel() error {
	db := global.GVA_MYSQL
	err := db.AutoMigrate(
		bot.BotBanContent{},
		bot.Bot{},
		bot.BanRecord{},
		bot.BotChatGroup{},
		bot.BotChatGroupClassify{},
		bot.BotBanGroupMem{},
		bot.BotTask{},
		bot.BotChannel{},
		bot.BotCmdConfig{},
		bot.BotChatGroupRelatedChannelFollow{},
		bot.BotMassMsgRecord{},
		bot.BotMassMsgPermission{},
		recharge.RechargeConfig{},
		recharge.UserRechargeRecord{},
		recharge.UserWallet{},
		recharge.AdPublishRecord{},
		ledger.Ledger{},
		ledger.LedgerAccountGroup{},
		ledger.LedgerPermission{},
		ledger2.Ledger{},
		ledger2.LedgerPermission{},
		ledger2.LedgerSession{},
	)
	if err != nil {
		return err
	}
	return nil
}
