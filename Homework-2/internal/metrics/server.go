package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen(host string, reg *prometheus.Registry) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	return http.ListenAndServe(host, mux)
}

func PickpointCounter() prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pickpoint_grpc",
		Help: "Number of requests handled",
	})
}

func GivenOrdersGauge() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "given_orders_grpc",
		Help: "Number of given orders",
	})
}

func RequestPickpointHistogram() prometheus.Histogram {
	return prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "pickpoint",
		Subsystem: "grpc",
		Name:      "request",
		Help:      "Requests handling histogram",
	})
}

func FailedOrderCounter() prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_orders_grpc",
		Help: "Number of failed requests to order service",
	})
}
