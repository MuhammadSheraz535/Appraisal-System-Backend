package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	logrus "github.com/sirupsen/logrus"
)

// const LoggerTimeStampFormat = "2006-01-02T15:04:05.999Z07:00"
const LoggerTimeStampFormat = time.RFC1123Z

type Fields map[string]interface{}

func TextLogInit() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat: LoggerTimeStampFormat,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
	logrus.SetFormatter(formatter)
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetLevel(getLogLevel())
}

func JSONLogInit() {
	logrus.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{
		TimestampFormat: LoggerTimeStampFormat,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return funcName, fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
	//logrus.SetLevel(logrus.TraceLevel)
	logrus.SetLevel(getLogLevel())
}

// getLogCaller
// setReportCaller cannot ignore caller levels, so cannot work in a wrapper. Fix not merged in github, so doing it manually
// func getLogCaller() string {
// 	_, file, line, _ := runtime.Caller(2)

// 	location := fmt.Sprintf("%v:%v", path.Base(file), line)
// 	return location
// }

func Info(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Info(args...)
}

func Debug(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Debug(args...)
}

func Trace(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Trace(args...)
}

func Warn(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Warn(args...)
}

func Error(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Error(args...)
}

func Panic(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Panic(args...)
}

func Fatal(args ...interface{}) {
	TextLogInit()
	logrus.WithFields(logrus.Fields{
		// "from": getLogCaller(),
	}).Fatal(args...)
}

func WithFields(fields Fields) *logrus.Entry {
	TextLogInit()
	return logrus.WithFields(logrus.Fields(fields))
}

func getLogLevel() logrus.Level {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
