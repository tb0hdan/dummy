package service_healthcheck

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/akhripko/dummy/log"
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

func New(port int, healthChecks []func() error, readinessChecks []func() error) Service {
	return &service{
		http: &http.Server{
			Addr:    fmt.Sprintf("0.0.0.0:%d", port),
			Handler: handler(healthChecks, readinessChecks),
		},
	}
}

func handler(healthChecks []func() error, readinessChecks []func() error) http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/version", serveVersion)
	handler.HandleFunc("/", func(res http.ResponseWriter, _ *http.Request) { serveCheck(res, healthChecks) })
	handler.HandleFunc("/ready", func(res http.ResponseWriter, _ *http.Request) { serveCheck(res, readinessChecks) })
	return handler
}

func serveVersion(response http.ResponseWriter, _ *http.Request) {
	writeFile("version", response)
}

func writeFile(file string, response http.ResponseWriter) {
	if proto, err := ioutil.ReadFile(file); err == nil { // nolint
		response.WriteHeader(http.StatusOK)
		response.Write(proto) // nolint
	} else {
		response.WriteHeader(http.StatusNoContent)
	}
}

func serveCheck(response http.ResponseWriter, checks []func() error) {
	writtenHeader := false
	for _, check := range checks {
		if err := check(); err != nil {
			if !writtenHeader {
				response.WriteHeader(http.StatusInternalServerError)
				writtenHeader = true
			}
			response.Write([]byte(err.Error())) // nolint
			response.Write([]byte("\n\n"))      // nolint
		}
	}

	if !writtenHeader {
		response.WriteHeader(http.StatusNoContent)
	}
}
