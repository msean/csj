package global

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func NewDB(dbconf DBConf) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&checkConnLiveness=true&loc=Local",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.DB,
	)
	logLevel := logger.Silent
	if dbconf.Debug {
		logLevel = logger.Info
	}

	var _logger logger.Writer
	_logger = GlobalRunTime.Logger

	// if GlobalRunTime.IsDebugMode() {
	// 	_logger = log.New(os.Stdout, "\r\n", log.LstdFlags)
	// }

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			_logger,
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,    // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		),
	})
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
