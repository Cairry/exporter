package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

// BaseDescCollector 统一管理多个指标的 Describe 和 Collect
type BaseDescCollector struct {
	metrics []prometheus.Metric
	descs   []*prometheus.Desc
}

func NewBaseDescCollector(registry *prometheus.Registry) *BaseDescCollector {
	base := &BaseDescCollector{
		metrics: make([]prometheus.Metric, 0),
		descs:   make([]*prometheus.Desc, 0),
	}

	registry.MustRegister(base)
	return base
}

// AddMetric 添加一个新的指标
func (b *BaseDescCollector) AddMetric(desc *prometheus.Desc, val float64, labelValues ...string) {
	b.metrics = append(b.metrics, prometheus.MustNewConstMetric(
		desc,
		prometheus.CounterValue,
		val,
		labelValues...,
	))
	b.descs = append(b.descs, desc)
}

// Describe 实现 prometheus.Collector 接口
func (b *BaseDescCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range b.descs {
		ch <- desc
	}
}

// Collect 实现 prometheus.Collector 接口
func (b *BaseDescCollector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range b.metrics {
		ch <- metric
	}
}
