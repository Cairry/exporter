package gauge

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicGauge struct {
	gaugeVec *prometheus.GaugeVec
}

func NewDynamicGauge(base *collectors.BaseCollector) *DynamicGauge {
	gaugeVec := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   "",
			Subsystem:   "",
			Name:        "example_dynamic_gauge",
			Help:        "A gauge with dynamic labels",
			ConstLabels: nil,
		},
		[]string{"job", "instance"}, // value 可变
	)

	dc := &DynamicGauge{
		gaugeVec: gaugeVec,
	}

	base.Register(dc.gaugeVec)
	return dc
}

func (dc *DynamicGauge) Inc(job, instance string) {
	dc.gaugeVec.WithLabelValues(job, instance).Inc()
}

func (dc *DynamicGauge) Dec(job, instance string) {
	dc.gaugeVec.WithLabelValues(job, instance).Dec()
}

func (dc *DynamicGauge) Set(job, instance string, value float64) {
	dc.gaugeVec.WithLabelValues(job, instance).Set(value)
}
