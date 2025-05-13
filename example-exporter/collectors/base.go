package collectors

import "github.com/prometheus/client_golang/prometheus"

// BaseCollector 是通用 Counter Collector，支持统一 Describe 和 Collect
type BaseCollector struct {
	registry   *prometheus.Registry
	collectors []prometheus.Collector
}

func NewBaseCollector(registry *prometheus.Registry) *BaseCollector {
	return &BaseCollector{
		registry:   registry,
		collectors: make([]prometheus.Collector, 0),
	}
}

// Register 注册任意 collector 到 Prometheus 中
func (b *BaseCollector) Register(c prometheus.Collector) {
	b.collectors = append(b.collectors, c)
	b.registry.MustRegister(c)
}

// Describe 实现 prometheus.Collector 接口
func (b *BaseCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range b.collectors {
		c.Describe(ch)
	}
}

// Collect 实现 prometheus.Collector 接口
func (b *BaseCollector) Collect(ch chan<- prometheus.Metric) {
	for _, c := range b.collectors {
		c.Collect(ch)
	}
}
