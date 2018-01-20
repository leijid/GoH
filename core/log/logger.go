package log

import (
	"GoH/core/file"
	"fmt"
	"github.com/op/go-logging"
	"os"
	"path"
)

var log = logging.MustGetLogger("goh")

var format = logging.MustStringFormatter(
	`%{color}%{time} %{shortfunc} > %{level:.4s} %{pid}%{color:reset} %{message}`,
)

func NewLogger(pathName string, level string) {
	dirExist, _ := file.PathExists(path.Dir(pathName))
	if !dirExist {
		os.MkdirAll(path.Dir(pathName), 0777)
	}
	os.Create(pathName)
	logFile, err := os.OpenFile(pathName, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	fileLog := logging.NewLogBackend(logFile, "", 0)
	stdLog := logging.NewLogBackend(os.Stderr, "", 0)
	stdFormatter := logging.NewBackendFormatter(stdLog, format)
	fileLeveled := logging.AddModuleLevel(fileLog)
	lLevel, _ := logging.LogLevel(level)
	fileLeveled.SetLevel(lLevel, "")
	logging.SetBackend(fileLeveled, stdFormatter)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Notice(args ...interface{}) {
	log.Notice(args...)
}

func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Warning(args ...interface{}) {
	log.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Critical(args ...interface{}) {
	log.Critical(args...)
}

func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

func CheckError(err error) {
	if err != nil {
		Errorf("Fatal error: %s", err.Error())
	}
}
