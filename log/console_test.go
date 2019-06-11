package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsoleError(t *testing.T) {
	var buf bytes.Buffer
	c := ConsoleLogger{
		lvl: ErrorLvl,
		w:   &buf,
	}
	c.Debug("debug")
	c.Info("info")
	c.Error("error")
	assert.Equal(t, "error\n", buf.String())
}

func TestConsoleInfo(t *testing.T) {
	var buf bytes.Buffer
	c := ConsoleLogger{
		lvl: InfoLvl,
		w:   &buf,
	}
	c.Debug("debug")
	c.Info("info")
	c.Error("error")
	assert.Equal(t, "info\nerror\n", buf.String())
}

func TestConsoleDebug(t *testing.T) {
	var buf bytes.Buffer
	c := ConsoleLogger{
		lvl: DebugLvl,
		w:   &buf,
	}
	c.Debug("debug")
	c.Info("info")
	c.Error("error")
	assert.Equal(t, "debug\ninfo\nerror\n", buf.String())
}
