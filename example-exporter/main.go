package main

import (
	"exporter/collectors"
	"exporter/collectors/counter"
	"exporter/collectors/gauge"
	"exporter/collectors/histogram"
	"exporter/collectors/summary"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	emptyRegister    = prometheus.NewRegistry()
	StaticCounter    *counter.StaticCounter
	DynamicCounter   *counter.DynamicCounter
	StaticGauge      *gauge.StaticGauge
	DynamicGauge     *gauge.DynamicGauge
	StaticHistogram  *histogram.StaticHistogram
	DynamicHistogram *histogram.DynamicHistogram
	StaticSummary    *summary.StaticSummary
	DynamicSummary   *summary.DynamicSummary
	DescCounter      *counter.DescCounter
	DescGauge        *gauge.DescGauge
)

func main() {
	fmt.Println("successfully")

	baseCollector := collectors.NewBaseCollector(emptyRegister)
	{
		{
			StaticCounter = counter.NewStaticCounter(baseCollector)
			DynamicCounter = counter.NewDynamicCounter(baseCollector)
		}
		{
			StaticGauge = gauge.NewStaticGauge(baseCollector)
			DynamicGauge = gauge.NewDynamicGauge(baseCollector)
		}
		{
			StaticHistogram = histogram.NewStaticHistogram(baseCollector)
			DynamicHistogram = histogram.NewDynamicHistogram(baseCollector)
		}
		{
			StaticSummary = summary.NewStaticSummary(baseCollector)
			DynamicSummary = summary.NewDynamicSummary(baseCollector)
		}
	}

	baseDescCollector := collectors.NewBaseDescCollector(emptyRegister)
	{
		{
			DescCounter = counter.NewDescCounter(baseDescCollector)
		}
		{
			DescGauge = gauge.NewDescGauge(baseDescCollector)
		}
	}

	http.HandleFunc("/metrics", handleMetrics)
	http.HandleFunc("/value", handleValue)

	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		return
	}

}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
	fmt.Println("--- 开始抓取 ---")
	promhttp.HandlerFor(emptyRegister,
		promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}

func handleValue(w http.ResponseWriter, r *http.Request) {
	StaticCounter.Inc()
	DynamicCounter.Inc("test1", "zux")
	StaticGauge.Inc()
	DynamicGauge.Inc("test2", "zux")
	StaticHistogram.Observe(1)
	DynamicHistogram.WithLabelValuesObserve("test3", 1)
	StaticSummary.Observe(1)
	DynamicSummary.WithLabelValuesObserve("test4", 1)
	DescCounter.Inc([]string{"test5"})
	DescGauge.Inc([]string{"test6"})
}
