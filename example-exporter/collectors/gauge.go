package collectors

import "github.com/prometheus/client_golang/prometheus"

type GaugeCollect struct {
	gaugeDesc    *prometheus.Desc
	gaugeMetrics prometheus.Metric
}

func NewGaugeCollect() *GaugeCollect {

	return &GaugeCollect{
		gaugeDesc: prometheus.NewDesc(
			"testGauge",
			"testGauge",
			[]string{"test"},
			nil,
		),
	}

}

func (g *GaugeCollect) Describe(docs chan<- *prometheus.Desc) {

	docs <- g.gaugeDesc

}

func (g *GaugeCollect) Collect(metrics chan<- prometheus.Metric) {

	metrics <- g.NewCollectMetrics()

}

func (g *GaugeCollect) NewCollectMetrics() prometheus.Metric {

	g.gaugeMetrics = prometheus.MustNewConstMetric(
		g.gaugeDesc,
		prometheus.GaugeValue,
		1,
		"test",
	)

	return g.gaugeMetrics

}
