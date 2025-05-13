package summary

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicSummary struct {
	summary *prometheus.SummaryVec
}

func NewDynamicSummary(base *collectors.BaseCollector) *DynamicSummary {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "example_dynamic_summary",
			Help:      "A summary with dynamic labels",
			ConstLabels: prometheus.Labels{
				"job": "dynamic_job",
			},
			Objectives: map[float64]float64{
				0.1: 0, 0.2: 0, 0.3: 0,
				0.4: 0, 0.5: 0, 0.6: 0,
				0.7: 0, 0.8: 0, 0.88: 0,
				0.9: 0, 0.95: 0},
			MaxAge:     0,
			AgeBuckets: 0,
			BufCap:     0,
		},
		[]string{"app"},
	)

	ds := &DynamicSummary{
		summary: summary,
	}

	base.Register(ds.summary)
	return ds
}

// WithObserve 向 bucket 中添加值，动态标签，适合从外部传入；
func (d *DynamicSummary) WithObserve(labels prometheus.Labels, value float64) {
	d.summary.With(labels).Observe(value)
}

// WithLabelValuesObserve 向 bucket 中添加值，固定顺序标签，适合简单封装；
func (d *DynamicSummary) WithLabelValuesObserve(app string, value float64) {
	d.summary.WithLabelValues(app).Observe(value)
}
