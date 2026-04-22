package system

import (
	"github.com/msean/csj/backend/global"
	"github.com/msean/csj/backend/model/system/response"
)

type AutoCodeService struct{}

type Database interface {
	GetDB(businessDB string) (data []response.Db, err error)
	GetTables(businessDB string, dbName string) (data []response.Table, err error)
	GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error)
}

func (autoCodeService *AutoCodeService) Database(businessDB string) Database {

	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return AutoCodeMysql
	case "pgsql":
		return AutoCodePgsql
	case "mssql":
		return AutoCodeMssql
	case "oracle":
		return AutoCodeOracle
	case "sqlite":
		return AutoCodeSqlite
	default:
		return AutoCodeMysql
	}
}
