package histogram

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type DynamicHistogram struct {
	histogram *prometheus.HistogramVec
}

func NewDynamicHistogram(base *collectors.BaseCollector) *DynamicHistogram {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "example_dynamic_histogram",
		Help:      "A histogram with dynamic labels",
		ConstLabels: prometheus.Labels{
			"job": "dynamic_job",
			"env": "production",
		},
		Buckets: []float64{
			0.05, 0.1, 0.25,
			0.5, 1, 2.1, 2.5,
			5, 5.4, 10,
		},
		NativeHistogramBucketFactor:     0,
		NativeHistogramZeroThreshold:    0,
		NativeHistogramMaxBucketNumber:  0,
		NativeHistogramMinResetDuration: 0,
		NativeHistogramMaxZeroThreshold: 0,
	},
		[]string{"app"},
	)

	dh := &DynamicHistogram{
		histogram: histogram,
	}

	base.Register(dh.histogram)
	return dh
}

// WithObserve 向 bucket 中添加值，动态标签，适合从外部传入；
func (d *DynamicHistogram) WithObserve(labels prometheus.Labels, value float64) {
	d.histogram.With(labels).Observe(value)
}

// WithLabelValuesObserve 向 bucket 中添加值，固定顺序标签，适合简单封装；
func (d *DynamicHistogram) WithLabelValuesObserve(app string, value float64) {
	d.histogram.WithLabelValues(app).Observe(value)
}
