package log

import (
	"fmt"
	"strings"
)

type Logger interface {
	Error(v ...interface{})
	Info(v ...interface{})
	Debug(v ...interface{})
}

type Level int

const (
	ErrorLvl Level = 1
	InfoLvl  Level = 2
	DebugLvl Level = 3
)

func StrToLogLevel(s string) Level {
	switch strings.ToLower(s) {
	case "error":
		return ErrorLvl
	case "info":
		return InfoLvl
	default:
		return DebugLvl
	}
}

func SetLogger(l Logger) {
	logger = l
}

var logger Logger

func Error(v ...interface{}) {
	if logger == nil {
		fmt.Println(v...)
		return
	}
	logger.Error(v...)
}

func Info(v ...interface{}) {
	if logger == nil {
		fmt.Println(v...)
		return
	}
	logger.Info(v...)
}

func Debug(v ...interface{}) {
	if logger == nil {
		fmt.Println(v...)
		return
	}
	logger.Debug(v...)
}
