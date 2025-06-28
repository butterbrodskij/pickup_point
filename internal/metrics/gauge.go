package metrics

import "github.com/prometheus/client_golang/prometheus"

type Gauge struct {
	prometheus.Gauge
}

type UnImplementedGauge struct {
}

func GivenOrdersGauge() Gauge {
	return Gauge{prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "given_orders_grpc",
		Help: "Number of given orders",
	})}
}

func (g *Gauge) SuccessGaugeAdd(err error, add float64) {
	if err != nil {
		g.Add(add)
	}
}

func (g *Gauge) SuccessGaugeDec(err error) {
	if err != nil {
		g.Dec()
	}
}

func (g *UnImplementedGauge) SuccessGaugeAdd(err error, add float64) {
}

func (g *UnImplementedGauge) SuccessGaugeDec(err error) {
}
