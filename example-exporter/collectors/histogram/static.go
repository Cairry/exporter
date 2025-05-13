package histogram

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type StaticHistogram struct {
	histogram prometheus.Histogram
}

func NewStaticHistogram(base *collectors.BaseCollector) *StaticHistogram {
	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "",                               // 指标命名空间（可选），用于组织指标名称层级
			Subsystem: "",                               // 子系统名，常用于进一步划分模块
			Name:      "example_static_histogram",       // 指标名（必填），需唯一且语义明确
			Help:      "A histogram with static labels", // 帮助信息，用于描述该指标用途
			ConstLabels: prometheus.Labels{
				"job": "static_job",
				"env": "production",
			}, // 固定标签（静态标签），适用于整个指标都共有的 label
			Buckets: []float64{
				0.05, 0.1, 0.25,
				0.5, 1, 2.1, 2.5,
				5, 5.4, 10,
			}, // 自定义 bucket 列表（非 Native Histogram 使用）
			NativeHistogramBucketFactor:     0, // Native Histogram 配置：桶增长因子（默认自动）
			NativeHistogramZeroThreshold:    0, // Native Histogram 中“零值”判定阈值
			NativeHistogramMaxBucketNumber:  0, // Native Histogram 最大桶数限制
			NativeHistogramMinResetDuration: 0, // Native Histogram 桶最小重置时间间隔
			NativeHistogramMaxZeroThreshold: 0, // Native Histogram 零值最大阈值
		},
	)

	sh := &StaticHistogram{
		histogram: histogram,
	}

	base.Register(sh.histogram)
	return sh
}

// Observe 向 bucket 内添加值
func (s *StaticHistogram) Observe(float64 float64) {
	s.histogram.Observe(float64)
}
