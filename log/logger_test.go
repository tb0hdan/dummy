package log

import (
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
