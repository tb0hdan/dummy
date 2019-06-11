package service_prometheus

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Service interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}

type service struct {
	http *http.Server
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
