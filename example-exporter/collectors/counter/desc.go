package counter

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DescCounter struct {
	base   *collectors.BaseDescCollector
	desc   *prometheus.Desc
	metric prometheus.Metric
	value  float64
}

func NewDescCounter(base *collectors.BaseDescCollector) *DescCounter {
	desc := prometheus.NewDesc(
		"example_desc_counter",
		"A counter with desc labels",
		[]string{"url"}, // label tag，可变值
		map[string]string{ // 固定的 label
			"code": "200",
		},
	)

	c := &DescCounter{
		base: base,
		desc: desc,
	}

	return c
}

// Inc 值 +1
func (d *DescCounter) Inc(labelValues []string) {
	d.value++
	d.base.AddMetric(d.desc, d.value, labelValues...)
}

// Add 值 +n
func (d *DescCounter) Add(labelValues []string, value float64) {
	d.value += value
	d.base.AddMetric(d.desc, d.value, labelValues...)
}
