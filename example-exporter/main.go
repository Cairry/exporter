package main

import (
	"exporter/collectors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	emptyRegister = prometheus.NewRegistry()
)

func main() {

	fmt.Println("successfully")

	emptyRegister.MustRegister(
		collectors.NewHistogramCollect(),
		collectors.NewCounterCollectDesc(),
		collectors.NewGaugeCollect(),
		collectors.NewSummaryCollect(),
		collectors.NewSummary2Collect(),
		collectors.NewHistogram2Collect(),
		collectors.NewGauge2Collect(),
		collectors.NewCounter2Collect(),
	)

	http.HandleFunc("/metrics", handleMetrics)

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
