package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCounts prometheus.Counter
	RequestTiming prometheus.Histogram
	SomeGauge     prometheus.Gauge
)

func Register() {
	RequestCounts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_counts",
		Help: "Number of requests",
	})

	RequestTiming = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "request_timing",
		Help:    "The request timing",
		Buckets: prometheus.ExponentialBuckets(0.5, 2, 15),
	})

	SomeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "some_gauge",
		Help: "The gauge for something",
	})

	prometheus.MustRegister(
		RequestTiming,
		RequestCounts,
		SomeGauge,
	)
}

func Unregister() { // nolint megacheck
	prometheus.Unregister(RequestTiming)
	prometheus.Unregister(RequestCounts)
	prometheus.Unregister(SomeGauge)
}
