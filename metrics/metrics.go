package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCounts                         prometheus.Counter
	TimeIntervalRequested                 prometheus.Histogram
	TimeIntervalTotalConsumptionRequested prometheus.Histogram
)

func Register() {
	RequestCounts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_counts",
		Help: "Number of requests",
	})

	TimeIntervalRequested = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "time_interval_requested_hours",
		Help:    "The time interval requested in hours",
		Buckets: prometheus.ExponentialBuckets(0.5, 2, 15),
	})

	TimeIntervalTotalConsumptionRequested = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "time_interval_total_consumption_requested_hours",
		Help:    "The time interval total consumption requested in hours",
		Buckets: prometheus.ExponentialBuckets(0.5, 2, 15),
	})

	prometheus.MustRegister(TimeIntervalRequested, TimeIntervalTotalConsumptionRequested, RequestCounts)
}

func Unregister() { // nolint megacheck
	prometheus.Unregister(TimeIntervalRequested)
	prometheus.Unregister(TimeIntervalTotalConsumptionRequested)
	prometheus.Unregister(RequestCounts)
}
