package gauge

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DescGauge struct {
	base   *collectors.BaseDescCollector
	desc   *prometheus.Desc
	metric prometheus.Metric
	value  float64
}

func NewDescGauge(base *collectors.BaseDescCollector) *DescGauge {
	desc := prometheus.NewDesc(
		"example_desc_gauge",
		"A gauge with desc labels",
		[]string{"app"},
		nil,
	)

	dg := &DescGauge{
		base: base,
		desc: desc,
	}

	return dg
}

// Inc 值 +1
func (d *DescGauge) Inc(labelValues []string) {
	d.value++
	d.base.AddMetric(d.desc, d.value, labelValues...)
}

// Dec 值 -1
func (d *DescGauge) Dec(labelValues []string) {
	d.value--
	d.base.AddMetric(d.desc, d.value, labelValues...)
}

// Add 值 +n
func (d *DescGauge) Add(labelValues []string, value float64) {
	d.value += value
	d.base.AddMetric(d.desc, d.value, labelValues...)
}

// Sub 值 -n
func (d *DescGauge) Sub(labelValues []string, value float64) {
	d.value -= value
	d.base.AddMetric(d.desc, d.value, labelValues...)
}
