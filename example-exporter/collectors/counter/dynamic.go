package counter

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicCounter struct {
	counterVec *prometheus.CounterVec
}

func NewDynamicCounter(base *collectors.BaseCollector) *DynamicCounter {
	counterVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "example_dynamic_counter",
			Help:        "A counter with dynamic labels",
			ConstLabels: nil,
		},
		[]string{"job", "instance"}, // value 可变
	)

	dc := &DynamicCounter{
		counterVec: counterVec,
	}

	base.Register(dc.counterVec)
	return dc
}

// Inc 值 +1
func (dc *DynamicCounter) Inc(job, instance string) {
	dc.counterVec.WithLabelValues(job, instance).Inc()
}

// Add 值 +n
func (dc *DynamicCounter) Add(job, instance string, value float64) {
	dc.counterVec.WithLabelValues(job, instance).Add(value)
}
