package collectors

import "github.com/prometheus/client_golang/prometheus"

type Histogram2Collect struct {
	histogram2Collect prometheus.Histogram
}

func NewHistogram2Collect() *Histogram2Collect {

	return &Histogram2Collect{
		histogram2Collect: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "testHistogram2",
			Help: "testHistogram2",
			Buckets: []float64{
				0.05, 0.1, 0.25,
				0.5, 1, 2.1, 2.5,
				5, 5.4, 10,
			},
		}),
	}

}

func (h2 *Histogram2Collect) Describe(docs chan<- *prometheus.Desc) {

	h2.histogram2Collect.Describe(docs)

}

func (h2 *Histogram2Collect) Collect(metrics chan<- prometheus.Metric) {

	h2.histogram2Collect.Collect(metrics)

	info := []float64{0.01, 0.02, 0.3, 0.4, 0.45, 1.5, 2.58, 4.9, 5.3, 11}
	for _, v := range info {
		h2.histogram2Collect.Observe(v)
	}

}
