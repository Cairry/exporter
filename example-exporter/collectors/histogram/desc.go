package histogram

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DescHistogram struct {
	desc   *prometheus.Desc
	metric prometheus.Metric
}

func NewHistogramCollect() *DescHistogram {
	return &DescHistogram{
		desc: prometheus.NewDesc(
			"testHistogram",
			"testHistogram",
			nil,
			nil,
		),
	}
}

func (c *DescHistogram) Describe(docs chan<- *prometheus.Desc) {
	docs <- c.desc
}

func (c *DescHistogram) Collect(metrics chan<- prometheus.Metric) {
	metrics <- c.NewCollectMetrics()
}

func (c *DescHistogram) NewCollectMetrics() prometheus.Metric {
	c.metric = prometheus.MustNewConstHistogram(
		c.desc,
		100, // 记录所有观测值的总和；用于 Prometheus 的 histogram_quantile() 和 rate() 等函数；如果你观测了 [100, 50, 30]，sum = 180。
		100, // 记录总共发生了多少次观测；用于计算分位数、平均值等；如果观测了 3 次，则 count = 3。
		map[float64]uint64{
			100: 1, // <=100ms 发生了 1 次
			200: 2, // >100ms && <=200ms 发生了 2 次
		},
	)

	return c.metric
}
