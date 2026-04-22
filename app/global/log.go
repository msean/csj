package global

import (
	"os"

	"github.com/robfig/cron"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapConf struct {
	Maxsize       int    `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"`                      // 单位 MB
	MaxBackups    int    `mapstructure:"max-backups" json:"max-backups" yaml:"max-backups"`          // 最大备份数
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                      // 最大保存天数
	Compress      bool   `mapstructure:"compress" json:"compress" yaml:"compress"`                   // 是否可以压缩
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Format        string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出
	FilePath      string `mapstructure:"file-path" json:"file-path"  yaml:"file-path"`               // 文件路径
	EncodeLevel   string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`       // 编码级
	StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"` // 栈名
	ShowLine      bool   `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                // 显示行
	LogInConsole  bool   `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"` // 输出控制台
	RetentionDay  int    `mapstructure:"retention-day" json:"retention-day" yaml:"retention-day"`    // 日志保留天数
}

// func (c *ZapConf) Encoder() zapcore.Encoder {
// 	config := zapcore.EncoderConfig{
// 		TimeKey:       "time",
// 		NameKey:       "name",
// 		LevelKey:      "level",
// 		CallerKey:     "caller",
// 		MessageKey:    "message",
// 		StacktraceKey: c.StacktraceKey,
// 		LineEnding:    zapcore.DefaultLineEnding,
// 		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
// 			encoder.AppendString(c.Prefix + t.Format("2006-01-02 15:04:05.000"))
// 		},
// 		EncodeLevel:    c.LevelEncoder(),
// 		EncodeCaller:   zapcore.FullCallerEncoder,
// 		EncodeDuration: zapcore.SecondsDurationEncoder,
// 	}
// 	if c.Format == "json" {
// 		return zapcore.NewJSONEncoder(config)
// 	}
// 	return zapcore.NewConsoleEncoder(config)

// }

// // LevelEncoder 根据 EncodeLevel 返回 zapcore.LevelEncoder
// func (c *ZapConf) LevelEncoder() zapcore.LevelEncoder {
// 	switch {
// 	case c.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
// 		return zapcore.LowercaseLevelEncoder
// 	case c.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
// 		return zapcore.LowercaseColorLevelEncoder
// 	case c.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
// 		return zapcore.CapitalLevelEncoder
// 	case c.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
// 		return zapcore.CapitalColorLevelEncoder
// 	default:
// 		return zapcore.LowercaseLevelEncoder
// 	}
// }

// func NewZapLogger(conf ZapConf) (logger *zap.Logger) {
// 	lumberjacklogger := &lumberjack.Logger{
// 		Filename:   conf.FilePath,
// 		MaxSize:    conf.Maxsize,
// 		MaxBackups: conf.MaxBackups,
// 		MaxAge:     conf.MaxAge,
// 		Compress:   conf.Compress, // disabled by default
// 	}
// 	defer lumberjacklogger.Close()

// 	fileEncoder := conf.Encoder()

// 	cores := []zapcore.Core{
// 		zapcore.NewCore(fileEncoder, zapcore.AddSync(lumberjacklogger), zap.InfoLevel),
// 	}
// 	if conf.LogInConsole {
// 		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
// 		consoleWriter := zapcore.Lock(os.Stdout)
// 		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel))
// 	}
// 	zapCores := zapcore.NewTee(cores...)

// 	var options []zap.Option
// 	if conf.ShowLine {
// 		options = append(options, zap.AddCaller())
// 	}
// 	logger = zap.New(zapCores, options...)
// 	defer logger.Sync()
// 	return
// }

func NewZapLogger(conf ZapConf) (logger *zap.Logger) {
	encoder := getEncoder()
	writeSyncer := getLogWriter(conf)
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

func getLogWriter(conf ZapConf) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   conf.FilePath,
		MaxSize:    conf.Maxsize,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
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

// 数据库日志
// var Xlogger *XormLogger

// type XormLogger struct {
// }

// func (this *XormLogger) Write(p []byte) (n int, err error) {
// 	Log.Info("数据库操作", zap.String("数据库", string(p)))
// 	return len(p), nil
// }

// // echo日志
// var EchoLog *EchoLogger

// type EchoLogger struct {
// }

// func (this *EchoLogger) Write(p []byte) (n int, err error) {
// 	Log.Info("ECHO", zap.String("请求", string(p)))
// 	return len(p), nil
// }
