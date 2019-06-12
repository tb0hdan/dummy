package service

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Timer(m prometheus.Observer, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	if m == nil {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		go m.Observe(time.Since(start).Seconds())
	}
}

func Counter(m prometheus.Counter, next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	if m == nil {
		return next
	}
	return func(w http.ResponseWriter, r *http.Request) {
		go m.Inc()
		next(w, r)
	}
}
