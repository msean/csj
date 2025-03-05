package global

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

// DBConf 数据库配置信息.
type DBConf struct {
	User     string `json:"user"`     // 链接用户名
	Password string `json:"password"` // 链接密码
	Host     string `json:"host"`     // 链接主机
	Port     string `json:"port"`     // 链接端口
	DB       string `json:"db"`       // 链接数据库
	Driver   string `json:"driver"`   // 链接类型
	Charset  string `json:"charset"`  // 字符集
	Debug    bool   `json:"debug"`    // 调试开关
}

func NewMysql(dbconf DBConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&checkConnLiveness=true&loc=Local",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.DB, dbconf.Charset,
	)
	logLevel := logger.Silent
	if dbconf.Debug {
		logLevel = logger.Info
	}

	var option gorm.Option
	Global.Logger.Info("Global.IsDebugMode()", zap.Bool("Global.IsDebugMode()", Global.IsDebugMode()))
	if Global.IsDebugMode() {
		option = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logLevel,    // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
					Colorful:                  true,        // Disable color
				},
			),
		}
	} else {
		l := zapgorm2.New(Global.Logger)
		l.LogLevel = logger.Info
		l.SetAsDefault()
		option = &gorm.Config{Logger: l}
	}

	db, err := gorm.Open(mysql.Open(dsn), option)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(500)
	sqlDB.SetConnMaxLifetime(time.Second * 300)
	return db, nil
}
