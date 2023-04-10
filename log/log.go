package log

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	xlog = logrus.New()
)

type Config struct {
	Dir   string
	Level string
}

func init() {
	Init(false)
}

func Init(debug bool) {
	InitLog(&Config{}, debug)
}

func InitLog(conf *Config, debug bool) {
	xlog.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	if conf.Dir != "" {
		pathMap := PathMap{
			logrus.InfoLevel:  conf.Dir + "/info.log",
			logrus.ErrorLevel: conf.Dir + "/error.log",
		}
		xlog.Hooks.Add(NewLocalHook(
			pathMap,
			&logrus.JSONFormatter{},
		))
	}
	if debug {
		xlog.SetLevel(logrus.DebugLevel)
	} else {
		if conf.Level != "" {
			slevel := strings.ToLower(conf.Level)
			switch slevel {
			case "fatal":
				xlog.SetLevel(logrus.FatalLevel)
			case "error":
				xlog.SetLevel(logrus.ErrorLevel)
			case "warn":
				xlog.SetLevel(logrus.WarnLevel)
			case "info":
				fallthrough
			default:
				xlog.SetLevel(logrus.InfoLevel)
			}
		} else {
			xlog.SetLevel(logrus.ErrorLevel)
		}
	}
}

func Print(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func Debug(format string, args ...interface{}) {
	xlog.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	xlog.Infof(format, args...)
}

func InfoUUID(uuid string, format string, args ...interface{}) {
	xlog.WithField("server_uuid", uuid).Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	xlog.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	xlog.WithField("stack", Stack(2, 1000)).Errorf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	xlog.WithField("stack", Stack(2, 1000)).Fatalf(format, args...)
}

func LevelEnabled(level logrus.Level) bool {
	return xlog.IsLevelEnabled(level)
}

func Entry() *logrus.Entry {
	return logrus.NewEntry(xlog)
}
