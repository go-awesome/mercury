//
//  logger.go
//  mercury
//
//  Copyright (c) 2016 Miguel Ángel Ortuño. All rights reserved.
//

package logger

import (
	"fmt"
	"os"
	"sync"
	"github.com/ivpusic/golog"
    "sync/atomic"
)

const loggerName = "github.com/ortuman/mercury"

const debugLevel = 0
const infoLevel  = 1
const warnLevel  = 2
const errorLevel = 3

type logger struct {
	logger  *golog.Logger
	level   int32
	log     chan struct{msg string; level int}
}

type logAppender struct {
	file *os.File
	bytesWritten int
}

// singleton interface

var instance *logger
var once sync.Once

const logQueueSize = 32768

func loggerInstance() *logger {
	once.Do(func() {
		instance = new(logger)
		instance.logger = golog.GetLogger(loggerName)
        instance.level = debugLevel
		instance.log = make(chan struct{msg string; level int}, logQueueSize)
		go instance.run()
	})
	return instance
}

func SetLogLevel(logLevel string) {
    l := loggerInstance()
	switch logLevel {
	case "DEBUG":
        atomic.StoreInt32(&l.level, debugLevel)
		loggerInstance().logger.Level = golog.DEBUG
	case "INFO":
        atomic.StoreInt32(&l.level, infoLevel)
        loggerInstance().logger.Level = golog.INFO
	case "WARN":
        atomic.StoreInt32(&l.level, warnLevel)
        loggerInstance().logger.Level = golog.WARN
	case "ERROR":
        atomic.StoreInt32(&l.level, errorLevel)
        loggerInstance().logger.Level = golog.ERROR
	}
}

func SetLogFilePath(logPath string) {
    loggerInstance().logger.Enable(newMercuryLogger(golog.Conf{
		// file in which logs will be saved
		"path": logPath,
	}))
}

func Debugf(msg string, params ...interface{}) {
    l := loggerInstance()
    if atomic.LoadInt32(&l.level) <= debugLevel {
        s := fmt.Sprintf(msg, params...)
        l.log <- struct{ msg string; level int }{msg: s, level: debugLevel}
    }
}

func Infof(msg string, params ...interface{}) {
    l := loggerInstance()
    if atomic.LoadInt32(&l.level) <= infoLevel {
        s := fmt.Sprintf(msg, params...)
        l.log <- struct{ msg string; level int }{msg: s, level: infoLevel}
    }
}

func Warnf(msg string, params ...interface{}) {
    l := loggerInstance()
    if atomic.LoadInt32(&l.level) <= warnLevel {
        s := fmt.Sprintf(msg, params...)
        l.log <- struct{ msg string; level int }{msg: s, level: warnLevel}
    }
}

func Errorf(msg string, params ...interface{}) {
    l := loggerInstance()
    if atomic.LoadInt32(&l.level) <= errorLevel {
        s := fmt.Sprintf(msg, params...)
        l.log <- struct{ msg string; level int }{msg: s, level: errorLevel}
    }
}

func (l *logger) run() {
	for log := range l.log {
		switch log.level {
		case debugLevel:
			l.logger.Debug(log.msg)
		case infoLevel:
			l.logger.Info(log.msg)
		case warnLevel:
			l.logger.Warn(log.msg)
		case errorLevel:
			l.logger.Error(log.msg)
		}
	}
}

func newMercuryLogger(cnf golog.Conf) *logAppender {
	path := cnf["path"]

	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &logAppender{file: f, bytesWritten: 0}
}

func (l *logAppender) Id() string {
	return "github.com/ortuman/mercury/logger"
}

func (l *logAppender) Append(log golog.Log) {
	if l.file == nil {
		return
	}

	var color string
	var icon string

	switch log.Level.Value {
	case golog.DEBUG.Value:
		color = "34m"
		icon = "★"
	case golog.INFO.Value:
		color = "32m"
		icon = "♥"
	case golog.WARN.Value:
		color = "33m"
		icon = "\u26A0"
	case golog.ERROR.Value:
		color = "31m"
		icon = "✖"
	default:
		return
	}

	msg := fmt.Sprintf("\033[36m%s\033[0m \033[37m%s\033[0m \033[%s%s[%s] ▶ %s\n\033[0m",
		log.Logger.Name,
		log.Time.Format("15:04:05"),
		color,
		icon,
		log.Level.Name[:4],
		log.Message)

	n, err := l.file.WriteString(msg)
	if err == nil {
		l.bytesWritten += n
		if l.bytesWritten >= 16384 {
			l.file.Sync()
			l.bytesWritten = 0
		}
	}
}
