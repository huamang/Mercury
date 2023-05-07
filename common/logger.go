package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	// 设置日志级别
	if Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
		PadLevelText:    true,
		DisableQuote:    true,
	})
	// 设置日志输出方式为标准输出
	logrus.SetOutput(os.Stdout)

}
