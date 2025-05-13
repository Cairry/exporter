package counter

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type StaticCounter struct {
	counter prometheus.Counter
}

func NewStaticCounter(base *collectors.BaseCollector) *StaticCounter {
	counter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "example_static_counter",
			Help:      "A counter with static labels",
			ConstLabels: prometheus.Labels{
				"job": "static_job",
				"env": "production",
			},
		},
	)

	sc := &StaticCounter{
		counter: counter,
	}

	base.Register(sc.counter)
	return sc
}

// Inc 值 +1
func (c *StaticCounter) Inc() {
	c.counter.Inc()
}

// Add 值 +n
func (c *StaticCounter) Add(value float64) {
	c.counter.Add(value)
}
