package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct {
	prometheus.Histogram
}

type UnImplementedHistogram struct {
}

func RequestPickpointHistogram() Histogram {
	return Histogram{prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "pickpoint",
		Subsystem: "grpc",
		Name:      "request",
		Help:      "Requests handling histogram",
	})}
}

func (h *Histogram) Observe(start time.Time) {
	h.Histogram.Observe(time.Since(start).Seconds())
}

func (h *UnImplementedHistogram) Observe(start time.Time) {
}
