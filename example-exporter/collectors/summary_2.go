package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Summary2Collect struct {
	summaryDuration prometheus.Summary
}

func NewSummary2Collect() *Summary2Collect {

	return &Summary2Collect{
		summaryDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Name: "test2Summary",
			Help: "test2Summary",
			Objectives: map[float64]float64{
				0.1: 0, 0.2: 0, 0.3: 0,
				0.4: 0, 0.5: 0, 0.6: 0,
				0.7: 0, 0.8: 0, 0.88: 0,
				0.9: 0, 0.95: 0},
		})}

}

func (s2 *Summary2Collect) Describe(docs chan<- *prometheus.Desc) {

	s2.summaryDuration.Describe(docs)

}

func (s2 *Summary2Collect) Collect(metrics chan<- prometheus.Metric) {

	s2.summaryDuration.Collect(metrics)

	info := []float64{0.4, 0.1, 0.6, 0.8, 1.0}
	for _, v := range info {
		s2.summaryDuration.Observe(v)
	}

}
