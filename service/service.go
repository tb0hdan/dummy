package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/metrics"

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
		Addr: fmt.Sprintf(":%d", port),
	}

	var srv service
	// initialize state
	go srv.initService()
	srv.setupHTTP(&httpSrv)

	return &srv
}

type service struct {
	http      *http.Server
	runErr    error
	readiness bool
}

func (s *service) setupHTTP(srv *http.Server) {
	srv.Handler = s.buildHandler()
	s.http = srv
}

func (s *service) buildHandler() http.Handler {
	r := mux.NewRouter()
	// path -> handlers

	// hello request
	hello := Counter(metrics.HelloRequestCounts, s.hello)
	hello = Timer(metrics.HelloRequestTiming, hello)
	r.HandleFunc("/hello", hello).Methods("GET")

	// ==============
	return r
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	log.Info("service: begin run")

	go func() {
		defer wg.Done()
		log.Debug("service addr:", s.http.Addr)
		err := s.http.ListenAndServe()
		if err != nil {
			s.runErr = err
			log.Error("service end run:", err)
			return
		}
		log.Info("service: end run")
	}()

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second) // nolint
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("service shutdown error:", err)
		}
	}()
}

func (s *service) HealthCheck() error {
	if !s.readiness {
		return errors.New("service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run service issue")
	}
	// TODO: add more checks like s.db.Ping()
	return nil
}

func (s *service) ReadinessCheck() error {
	if !s.readiness {
		return errors.New("service is't ready yet")
	}
	if s.runErr != nil {
		return errors.New("run service issue")
	}
	return nil
}
