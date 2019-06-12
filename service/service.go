package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
	"github.com/gorilla/mux"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	HealthCheck() error
	ReadinessCheck() error
}

func New(port int) Service {
	httpSrv := http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
	}

	var srv service
	// initialize state
	go srv.initService()
	srv.setupHttp(&httpSrv)

	return &srv
}

type service struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func (s *service) setupHttp(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *service) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers
	r.HandleFunc("/hello", s.hello).Methods("GET")
	// ==============
	return r
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("service: begin run")

	go func() {
		defer wg.Done()
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("service run error:", err)
			return
		}
		log.Info("service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("service shutdown error:", err)
		}
	}()
}

func (s *service) HealthCheck() error {
	if s.runErr != nil {
		return errors.New("run service issue")
	}
	return nil
}

func (s *service) ReadinessCheck() error {
	if s.runErr != nil {
		return errors.New("run service issue")
	}
	if s.readiness == false {
		return errors.New("service is't ready yet")
	}
	return nil
}
