package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type HistogramCollect struct {
	histogramDesc    *prometheus.Desc
	histogramMetrics prometheus.Metric
}

func NewHistogramCollect() *HistogramCollect {

	return &HistogramCollect{
		histogramDesc: prometheus.NewDesc(
			"testHistogram",
			"testHistogram",
			nil,
			nil,
		),
	}

}

func (c *HistogramCollect) Describe(docs chan<- *prometheus.Desc) {

	docs <- c.histogramDesc

}

func (c *HistogramCollect) Collect(metrics chan<- prometheus.Metric) {

	metrics <- c.NewCollectMetrics()

}

func (c *HistogramCollect) NewCollectMetrics() prometheus.Metric {

	c.histogramMetrics = prometheus.MustNewConstHistogram(
		c.histogramDesc,
		100,
		100,
		map[float64]uint64{100: 1, 200: 2},
	)

	return c.histogramMetrics

}
