package log

import (
	"github.com/akhripko/dummy/log/mock_log"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToLogLevel(t *testing.T) {
	assert.Equal(t, ErrorLvl, StrToLogLevel("Error"))
	assert.Equal(t, ErrorLvl, StrToLogLevel("ErroR"))
	assert.Equal(t, ErrorLvl, StrToLogLevel("ERROR"))
	assert.Equal(t, InfoLvl, StrToLogLevel("info"))
	assert.Equal(t, InfoLvl, StrToLogLevel("infO"))
	assert.Equal(t, InfoLvl, StrToLogLevel("INFO"))
	assert.Equal(t, DebugLvl, StrToLogLevel("debug"))
	assert.Equal(t, DebugLvl, StrToLogLevel("blabla"))
}

func TestError(t *testing.T) {
	c := gomock.NewController(t)
	l := mock_log.NewMockLogger(c)
	l.EXPECT().Error("data").Do(func(v interface{}) {
		assert.Equal(t, "data", v)
	}).Times(1)
	l.EXPECT().Debug("data").Times(0)
	l.EXPECT().Info("data").Times(0)
	SetLogger(l)
	Error("data")
}

func TestInfo(t *testing.T) {
	c := gomock.NewController(t)
	l := mock_log.NewMockLogger(c)
	l.EXPECT().Info("data").Do(func(v interface{}) {
		assert.Equal(t, "data", v)
	}).Times(1)
	l.EXPECT().Debug("data").Times(0)
	l.EXPECT().Error("data").Times(0)
	SetLogger(l)
	Info("data")
}

func TestDebug(t *testing.T) {
	c := gomock.NewController(t)
	l := mock_log.NewMockLogger(c)
	l.EXPECT().Debug("data").Do(func(v interface{}) {
		assert.Equal(t, "data", v)
	}).Times(1)
	l.EXPECT().Info("data").Times(0)
	l.EXPECT().Error("data").Times(0)
	SetLogger(l)
	Debug("data")
}
