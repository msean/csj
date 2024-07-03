package global

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func setLogger(level, filePath, filename string, maxAge, rotate time.Duration) *logrus.Logger {
	logger := logrus.New()
	writer, _ := rotatelogs.New(
		filePath,
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(maxAge),
		rotatelogs.WithRotationTime(rotate),
	)
	logger.SetOutput(writer)
	switch level {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}
