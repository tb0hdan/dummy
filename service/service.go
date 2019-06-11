package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
	"github.com/gorilla/mux"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	ReadRunError() error
}

func New(port int) Service {
	srv := service{
		http: &http.Server{
			Addr: fmt.Sprintf("0.0.0.0:%d", port),
		},
	}
	srv.initHandler()
	return &srv
}

type service struct {
	http   *http.Server
	runErr error
}

func (s *service) initHandler() {
	r := mux.NewRouter()
	// TODO: add rules
	s.http.Handler = r
}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := s.http.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
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

func (s *service) ReadRunError() error {
	return s.runErr
}
