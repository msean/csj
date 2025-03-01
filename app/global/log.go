package global

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Path string `mapstructure:"path"`
	Name string `mapstructure:"name"`
}

func GetLogFilePath(logConfig LogConfig) string {
	return fmt.Sprintf("%s/%s", logConfig.Path, logConfig.Name)
}

func InitLog(filepath string) *zap.Logger {
	encoder := getEncoder()
	writeSyncer := getLogWriter(filepath)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriter := zapcore.Lock(os.Stdout)
	//consoleDebug := zapcore.Lock(os.Stdout)
	//consoleEncodeer := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	//p := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
	//	return level >= zapcore.DebugLevel
	//})
	var allcode []zapcore.Core
	allcode = append(allcode, core)
	allcode = append(allcode, zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel))
	//allcode = append(allcode, zapcore.NewCore(consoleEncodeer, consoleDebug, p))
	c := zapcore.NewTee(allcode...)
	//zap.AddCaller() //添加将调用函数信息记录到日志中的功能。
	return zap.New(c, zap.AddCaller())
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filepath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    10240,
		MaxAge:     7,
		Compress:   true,
		LocalTime:  true,
		MaxBackups: 4,
	}

	c := cron.New()
	c.AddFunc("0 0 0 1/1 * ?", func() {
		lumberJackLogger.Rotate()
	})
	c.Start()
	return zapcore.AddSync(lumberJackLogger)
}
