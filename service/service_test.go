package service

import (
	"github.com/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_ReadinessCheck(t *testing.T) {
	var srv service

	srv.readiness = false
	assert.Equal(t, "service is't ready yet", srv.ReadinessCheck().Error())

	srv.readiness = true
	assert.Nil(t, srv.ReadinessCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "run service issue", srv.HealthCheck().Error())
}

func TestService_HealthCheck(t *testing.T) {
	var srv service

	srv.readiness = false
	assert.Equal(t, "service is't ready yet", srv.HealthCheck().Error())

	srv.readiness = true
	assert.Nil(t, srv.HealthCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "run service issue", srv.HealthCheck().Error())
}
