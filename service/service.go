package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	ReadRunError() error
}

// TODO: fix it
func New(port int) Service {
	srv := service{
		http: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
			Handler: handler(),
		},
	}
}

type service struct {
	http   *http.Server
	runErr error
}

func (s *service) initHandler() {

}

func (s *service) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		<-ctx.Done()
		sdCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := s.http.Shutdown(sdCtx)
		if err != nil {
			log.Error("service shutdown error:", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := s.http.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
	}()
}

func (s *service) ReadRunError() error {
	return s.runErr
}
