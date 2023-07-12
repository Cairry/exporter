package collectors

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	info sync.Map
)

type CounterCollect struct {
	counter200Desc *prometheus.Desc
	counter4xxDesc *prometheus.Desc
	counter5xxDesc *prometheus.Desc
	counterMetrics prometheus.Metric
}

func NewCounterCollectDesc() *CounterCollect {

	cc := &CounterCollect{
		counter200Desc: prometheus.NewDesc(
			"testCounter",
			"testCounter",
			[]string{"url"},
			map[string]string{"code": "200"},
		),
		counter4xxDesc: prometheus.NewDesc(
			"testCounter",
			"testCounter",
			[]string{"url"},
			map[string]string{"code": "4xx"},
		),
		counter5xxDesc: prometheus.NewDesc(
			"testCounter",
			"testCounter",
			[]string{"url"},
			map[string]string{"code": "5xx"},
		),
	}

	return cc

}

func (c *CounterCollect) Describe(docs chan<- *prometheus.Desc) {

	docs <- c.counter200Desc
	docs <- c.counter4xxDesc
	docs <- c.counter5xxDesc

}

func (c *CounterCollect) Collect(metrics chan<- prometheus.Metric) {

	state := getHTTPState()

	for _, key := range []string{"200", "4xx", "5xx"} {
		_, ok := info.Load(key)
		if !ok {
			info.Store(key, 0)
		}
	}

	switch string(strconv.Itoa(state)[0]) {
	case "2":
		value, _ := info.Load("200")
		newV := value.(int) + 1
		info.Store("200", newV)
	case "4":
		value, _ := info.Load("4xx")
		newV := value.(int) + 1
		info.Store("4xx", newV)
	case "5":
		value, _ := info.Load("5xx")
		newV := value.(int) + 1
		info.Store("5xx", newV)
	}

	num2xx := getMapNum("200")
	num4xx := getMapNum("4xx")
	num5xx := getMapNum("5xx")

	metrics <- c.NewCounterCollectMetrics(c.counter200Desc, num2xx)
	metrics <- c.NewCounterCollectMetrics(c.counter4xxDesc, num4xx)
	metrics <- c.NewCounterCollectMetrics(c.counter5xxDesc, num5xx)

}

func (c *CounterCollect) NewCounterCollectMetrics(docs *prometheus.Desc, num float64) prometheus.Metric {

	c.counterMetrics = prometheus.MustNewConstMetric(
		docs,
		prometheus.CounterValue,
		num,
		"http://url",
	)

	return c.counterMetrics
}

func getMapNum(key any) float64 {

	value, _ := info.Load(key)
	floatValue, _ := strconv.ParseFloat(strconv.Itoa(value.(int)), 64)

	return floatValue
}

func getHTTPState() int {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       20 * time.Second,
	}

	requestData := fmt.Sprintf(`{"messages": [{"role": "user", "content": "hello"}], "model": "%s", "max_tokens": 2048, "temperature": 0.7}`, "test")
	data := []byte(requestData)
	req, err := http.NewRequest("Post", "https://xx.js.design/", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return 0
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0
	}

	// 关闭请求
	defer resp.Body.Close()

	return resp.StatusCode

}
