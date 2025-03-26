package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"log2metric/storage"
	"log2metric/types"
)

type (
	LogLevelExporter struct {
		LabelStorage    *storage.CacheStorage
		LogGaugeCollect *prometheus.Desc
	}
)

func NewProcessGaugeCollect(labelStorage *storage.CacheStorage) *LogLevelExporter {
	return &LogLevelExporter{
		LabelStorage: labelStorage,
		LogGaugeCollect: prometheus.NewDesc(
			"l2m_level_info",
			"log level",
			[]string{"namespace", "service", "level"},
			nil,
		),
	}
}

func (l LogLevelExporter) Describe(docs chan<- *prometheus.Desc) {
	docs <- l.LogGaugeCollect
}

func (l LogLevelExporter) Collect(metric chan<- prometheus.Metric) {
	for _, label := range l.LabelStorage.List() {
		l.register(metric, label)
	}
}

func (l LogLevelExporter) register(metric chan<- prometheus.Metric, label types.MetricLabel) {
	metric <- prometheus.MustNewConstMetric(
		l.LogGaugeCollect,
		prometheus.GaugeValue,
		label.Value,
		label.Namespace,
		label.Service,
		label.Level,
	)
}
