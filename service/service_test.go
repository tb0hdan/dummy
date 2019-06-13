package service

import (
	"testing"

	"github.com/akhripko/dummy/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheck(t *testing.T) {
	var srv service

	srv.readiness = false
	assert.Equal(t, "service is't ready yet", srv.HealthCheck().Error())

	srv.readiness = true
	assert.Equal(t, "storage issue", srv.HealthCheck().Error())

	c := gomock.NewController(t)
	db := mock_service.NewMockStorage(c)
	db.EXPECT().Ping().DoAndReturn(func() error { return errors.New("some db error") }).Times(1)
	srv.storage = db
	assert.Equal(t, "storage issue", srv.HealthCheck().Error())

	c = gomock.NewController(t)
	db = mock_service.NewMockStorage(c)
	db.EXPECT().Ping().DoAndReturn(func() error { return nil }).Times(1)
	srv.storage = db
	assert.Nil(t, srv.HealthCheck())

	srv.runErr = errors.New("some run error")
	assert.Equal(t, "run service issue", srv.HealthCheck().Error())
}
