package service_prometheus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	HealthCheck() error
	ReadinessCheck() error
}

type service struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("prometheus service: begin run")

	go func() {
		defer wg.Done()
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("prometheus service run error:", err.Error())
			return
		}
		log.Info("prometheus service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("prometheus service shutdown error:", err.Error())
		}
	}()

	s.readiness = true
}

func New(port int) Service {
	return &service{
		http: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
			Handler: handler(),
		},
	}
}

func handler() http.Handler {
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	return handler
}

func (s *service) HealthCheck() error {
	if s.runErr != nil {
		return errors.New("run prometheus service issue")
	}
	return nil
}

func (s *service) ReadinessCheck() error {
	if s.runErr != nil {
		return errors.New("run prometheus service issue")
	}
	if s.readiness == false {
		return errors.New("prometheus service is't ready yet")
	}
	return nil
}
