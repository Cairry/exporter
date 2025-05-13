package summary

import (
	"exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
)

type StaticSummary struct {
	summary prometheus.Summary
}

func NewStaticSummary(base *collectors.BaseCollector) *StaticSummary {
	summary := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "",
			Subsystem: "",
			Name:      "example_static_summary",
			Help:      "A summary with static labels",
			ConstLabels: prometheus.Labels{
				"job": "static_job",
			},
			Objectives: map[float64]float64{
				0.1: 0, 0.2: 0, 0.3: 0,
				0.4: 0, 0.5: 0, 0.6: 0,
				0.7: 0, 0.8: 0, 0.88: 0,
				0.9: 0, 0.95: 0}, // 分位数目标配置，定义了分位数估计值及其各自的绝对误差值。Objectives[q] = e，则报告的 q 值将是 q-e 和 q+e 之间某个 φ 的 φ 分位数值。
			MaxAge:     0, // 数据保留时间（滑动窗口长度）
			AgeBuckets: 0, // 滑动窗口桶数量（与 MaxAge 配合使用）
			BufCap:     0, // 内部缓冲区容量
		},
	)

	ds := &StaticSummary{
		summary: summary,
	}

	base.Register(ds.summary)
	return ds
}

// Observe 向 bucket 内添加值
func (s *StaticSummary) Observe(value float64) {
	s.summary.Observe(value)
}
