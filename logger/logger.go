package logger

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	// Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.InfoLevel)
	Logger.Level = logrus.TraceLevel

	// src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	// if err != nil {
	// 	fmt.Printf("err=%v, prober aborted\n", err)
	// 	os.Exit(2)
	// }
	// logrus.SetOutput(src)
}
func LogToFile(fileName string) {

	// fileName := "/home/work/kafka-prober.log"

	writer, err := rotatelogs.New(
		fileName+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(60*60*24*7)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(1*60*60*24)*time.Second),
	)
	if err != nil {
		Logger.Fatal("Init log failed, err:", err)
	}
	Logger.SetOutput(writer)
}
