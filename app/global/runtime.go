package global

import (
	"app/pkg"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	EnvDebugMode = "debug"
	EnvTestMode  = "test"
)

var Global *RunTime

type RunTime struct {
	DB     *gorm.DB
	Viper  *viper.Viper
	Engine *gin.Engine
	server *http.Server
	Logger *zap.Logger
	Redis  redis.UniversalClient
	Sms    pkg.SmsSender
}

func WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		println("Config file changed:", e.Name)
	})
}

func InitViper(configFile string) (v *viper.Viper, err error) {
	if configFile == "" {
		err = errors.New("confilgFile not config")
		return
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("config")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	v = viper.GetViper()
	return
}

func InitRunTime(configfile string) (err error) {
	v, e := InitViper(configfile)
	if err != nil {
		return e
	}
	Global = &RunTime{
		Viper:  v,
		Engine: gin.Default(),
	}
	// logger 初始化
	var zapConf LogConfig
	if err = Global.applyConf("zap", &zapConf); err != nil {
		return
	}
	Global.Logger = InitLog(GetLogFilePath(zapConf))

	// mysql 初始化
	dbconf := DBConf{
		User:     Global.Viper.GetString("database.username"),
		Password: Global.Viper.GetString("database.password"),
		Host:     Global.Viper.GetString("database.host"),
		Port:     Global.Viper.GetString("database.port"),
		DB:       Global.Viper.GetString("database.db"),     // 链接数据库
		Driver:   Global.Viper.GetString("database.driver"), // 链接数据库,
		Charset:  Global.Viper.GetString("database.charset"),
		Debug:    Global.Viper.GetBool("database.debug"),
	}
	Global.DB, err = NewMysql(dbconf)

	// redis 初始化
	// GlobalRunTime.Redis = NewRedis(RedisConf{
	// 	Addr:     GlobalRunTime.Viper.GetString("redis.addr"),
	// 	Password: GlobalRunTime.Viper.GetString("redis.password"), // 密码
	// 	DB:       1,
	// })

	// sms 初始化
	Global.Sms = &pkg.AliPlatfrom{
		Password: Global.Viper.GetString("sms.password"),
		Uid:      Global.Viper.GetString("sms.uid"),
		Secret:   Global.Viper.GetString("sms.secret"),
	}
	return
}

func (r *RunTime) Signkey() string {
	return r.Viper.GetString("env.signkey")
}

func (r *RunTime) Env() string {
	return r.Viper.GetString("env.env")
}

func (r *RunTime) IsDebugMode() bool {
	return r.Env() == "debug"
}

func (r *RunTime) VerifyCode() string {
	return r.Viper.GetString("env.verifycode")
}

func (r *RunTime) Migrate() bool {
	return r.Viper.GetBool("database.migrate")
}

func (r *RunTime) TokenEx() time.Duration {
	return time.Minute * time.Duration(r.Viper.GetInt("env.tokenex"))
}

func (r *RunTime) Node() int {
	return r.Viper.GetInt("env.machineID")
}

func (r *RunTime) TokenExRefresh() time.Duration {
	return time.Minute * time.Duration(r.Viper.GetInt("env.tokenexrefresh"))
}

func (r *RunTime) LogFilepath() string {
	return r.Viper.GetString("logs.filepath")
}

func (r *RunTime) LogFileName() string {
	return r.Viper.GetString("logs.filename")
}

func (r *RunTime) SmsRegisterTemp() string {
	return r.Viper.GetString("sms.registerTemplateCode")
}

func (r *RunTime) SmsLoginTemp() string {
	return r.Viper.GetString("sms.loginTemplateCode")
}

func (r *RunTime) LogMaxAge() time.Duration {
	return time.Hour * time.Duration(r.Viper.GetInt("logs.maxage"))
}

func (r *RunTime) LogRotate() time.Duration {
	return time.Hour * time.Duration(r.Viper.GetInt("logs.rotate"))
}

func (r *RunTime) Run(router http.Handler) {
	r.server = &http.Server{
		Addr:    r.Viper.GetString("env.addr"),
		Handler: router,
	}
	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		r.Logger.DPanic(fmt.Sprintf("Server Run err: %s", err))
	}
}

func (r *RunTime) Close() {
	if err := r.server.Shutdown(context.Background()); err != nil {
		r.Logger.DPanic(fmt.Sprintf("Server Run err: %s", err))
	}
}

func (r *RunTime) applyConf(sub string, conf interface{}) (err error) {
	return r.Viper.Sub(sub).Unmarshal(&conf)
}
