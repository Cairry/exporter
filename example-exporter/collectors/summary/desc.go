package summary

import "github.com/prometheus/client_golang/prometheus"

type DescSummary struct {
	desc   *prometheus.Desc
	metric prometheus.Metric
}

func NewSummaryCollect() *DescSummary {
	return &DescSummary{
		desc: prometheus.NewDesc(
			"testSummary",
			"testSummary",
			[]string{"job"},
			nil,
		),
	}
}

func (s *DescSummary) Describe(docs chan<- *prometheus.Desc) {
	docs <- s.desc
}

func (s *DescSummary) Collect(metrics chan<- prometheus.Metric) {
	metrics <- s.NewCollectMetrics()
}

func (s *DescSummary) NewCollectMetrics() prometheus.Metric {
	s.metric = prometheus.MustNewConstSummary(
		s.desc,
		11,
		8,
		map[float64]float64{
			0.5:  0.01,
			0.9:  0.01,
			0.99: 0.001,
		},
		"test",
	)

	return s.metric
}
