package log

import (
	"fmt"
	"io"
	"os"
)

type ConsoleLogger struct {
	lvl Level
	w   io.Writer
}

func NewConsoleLogger(lvl Level) (Logger, error) {
	return &ConsoleLogger{
		lvl: lvl,
		w:   os.Stdout,
	}, nil
}

func (l *ConsoleLogger) Error(v ...interface{}) {
	if l.lvl >= ErrorLvl {
		fmt.Fprintln(l.w, v...)
	}
}

func (l *ConsoleLogger) Info(v ...interface{}) {
	if l.lvl >= InfoLvl {
		fmt.Fprintln(l.w, v...)
	}
}

func (l *ConsoleLogger) Debug(v ...interface{}) {
	if l.lvl >= DebugLvl {
		fmt.Fprintln(l.w, v...)
	}
}
