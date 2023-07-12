package collectors

import "github.com/prometheus/client_golang/prometheus"

type Counter2Collect struct {
	counter2Collect prometheus.Counter
}

func NewCounter2Collect() *Counter2Collect {

	return &Counter2Collect{
		counter2Collect: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "test2Counter",
			Help: "test2Counter",
		}),
	}

}

func (c2 *Counter2Collect) Describe(docs chan<- *prometheus.Desc) {

	c2.counter2Collect.Describe(docs)

}

func (c2 *Counter2Collect) Collect(metrics chan<- prometheus.Metric) {

	c2.counter2Collect.Collect(metrics)

	c2.counter2Collect.Inc()
	c2.counter2Collect.Add(2)

}
