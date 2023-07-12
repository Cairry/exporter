package collectors

import "github.com/prometheus/client_golang/prometheus"

type SummaryCollect struct {
	summaryDesc    *prometheus.Desc
	summaryMetrics prometheus.Metric
}

func NewSummaryCollect() *SummaryCollect {

	return &SummaryCollect{
		summaryDesc: prometheus.NewDesc(
			"testSummary",
			"testSummary",
			[]string{"test"},
			nil,
		),
	}

}

func (s *SummaryCollect) Describe(docs chan<- *prometheus.Desc) {

	docs <- s.summaryDesc

}

func (s *SummaryCollect) Collect(metrics chan<- prometheus.Metric) {

	metrics <- s.NewCollectMetrics()

}

func (s *SummaryCollect) NewCollectMetrics() prometheus.Metric {

	s.summaryMetrics = prometheus.MustNewConstSummary(
		s.summaryDesc,
		11,
		8,
		map[float64]float64{
			0.5:  0.01,
			0.9:  0.01,
			0.99: 0.001,
		},
		"test",
	)

	return s.summaryMetrics

}