package gauge

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type StaticGauge struct {
	gauge prometheus.Gauge
}

func NewStaticGauge(base *collectors.BaseCollector) *StaticGauge {
	gauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "example_static_gauge",
			Help:      "A Gauge with static labels",
			ConstLabels: prometheus.Labels{
				"job": "static_job",
				"env": "production",
			},
		},
	)

	sc := &StaticGauge{
		gauge: gauge,
	}

	base.Register(sc.gauge)
	return sc
}

// Set 设置规定值
func (s *StaticGauge) Set(value float64) {
	s.gauge.Set(value)
}

// Inc 值 +1
func (s *StaticGauge) Inc() {
	s.gauge.Inc()
}

// Dec 值 -1
func (s *StaticGauge) Dec() {
	s.gauge.Dec()
}

// Add 值 +n
func (s *StaticGauge) Add(value float64) {
	s.gauge.Add(value)
}

// Sub 值 -n
func (s *StaticGauge) Sub(value float64) {
	s.gauge.Sub(value)
}
