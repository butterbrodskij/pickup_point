package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Counter struct {
	prometheus.Counter
}

type UnImplementedCounter struct {
}

func PickpointCounter() Counter {
	return Counter{prometheus.NewCounter(prometheus.CounterOpts{
		Name: "pickpoint_grpc",
		Help: "Number of requests handled",
	})}
}

func FailedOrderCounter() Counter {
	return Counter{prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_orders_grpc",
		Help: "Number of failed requests to order service",
	})}
}

func (c *UnImplementedCounter) FailedCounterInc(err error) {
}

func (c *UnImplementedCounter) CounterInc() {
}

func (c *Counter) FailedCounterInc(err error) {
	if err != nil {
		c.Inc()
	}
}

func (c *Counter) CounterInc() {
	c.Inc()
}
