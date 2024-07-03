package global

import (
	"errors"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

const (
	EnvDebugMode = "debug"
	EnvTestMode  = "test"
)

var Conf *viper.Viper
var GlobalRunTime *RunTime

type RunTime struct {
	DB     *gorm.DB
	Viper  *viper.Viper
	Engine *gin.Engine
	Logger *logrus.Logger
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
	GlobalRunTime = &RunTime{
		Viper:  v,
		Engine: gin.Default(),
	}
	if !GlobalRunTime.IsDebugMode() {
		GlobalRunTime.Logger = setLogger(GlobalRunTime.Viper.GetString("logs.level"),
			GlobalRunTime.LogFilepath(), GlobalRunTime.LogFileName(), GlobalRunTime.LogMaxAge(), GlobalRunTime.LogRotate())
	} else {
		GlobalRunTime.Logger = logrus.New()
	}
	dbconf := DBConf{
		User:     GlobalRunTime.Viper.GetString("database.username"),
		Password: GlobalRunTime.Viper.GetString("database.password"),
		Host:     GlobalRunTime.Viper.GetString("database.host"),
		Port:     GlobalRunTime.Viper.GetString("database.port"),
		DB:       GlobalRunTime.Viper.GetString("database.db"),     // 链接数据库
		Driver:   GlobalRunTime.Viper.GetString("database.driver"), // 链接数据库,
		Charset:  GlobalRunTime.Viper.GetString("database.charset"),
		Debug:    GlobalRunTime.Viper.GetBool("database.debug"),
	}
	GlobalRunTime.DB, err = NewDB(dbconf)
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

func (r *RunTime) TokenExRefresh() time.Duration {
	return time.Minute * time.Duration(r.Viper.GetInt("env.tokenexrefresh"))
}

func (r *RunTime) LogFilepath() string {
	return r.Viper.GetString("logs.filepath")
}

func (r *RunTime) LogFileName() string {
	return r.Viper.GetString("logs.filename")
}

func (r *RunTime) LogMaxAge() time.Duration {
	return time.Hour * time.Duration(r.Viper.GetInt("logs.maxage"))
}

func (r *RunTime) LogRotate() time.Duration {
	return time.Hour * time.Duration(r.Viper.GetInt("logs.rotate"))
}

func (r *RunTime) Run() {
	r.Engine.Run(r.Viper.GetString("env.addr"))
}

func (r *RunTime) Close() {
}
