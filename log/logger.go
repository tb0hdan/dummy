package log

import "strings"

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
