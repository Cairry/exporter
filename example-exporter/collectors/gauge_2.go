package collectors

import "github.com/prometheus/client_golang/prometheus"

type Gauge2Collect struct {
	gauge2Collect prometheus.Gauge
}

func NewGauge2Collect() *Gauge2Collect {

	return &Gauge2Collect{
		gauge2Collect: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "test2Gauge",
			Help: "test2Gauge",
		}),
	}

}

func (g2 *Gauge2Collect) Describe(docs chan<- *prometheus.Desc) {

	g2.gauge2Collect.Describe(docs)

}

func (g2 *Gauge2Collect) Collect(metrics chan<- prometheus.Metric) {

	g2.gauge2Collect.Collect(metrics)

	g2.gauge2Collect.Inc()

}