package main

import (
	"Gin/exporter/config"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	// UrlStateCode url状态码
	UrlStateCode int
	// EmptyRegistry 清空默认指标
	EmptyRegistry = prometheus.NewRegistry()
)

// 带动态标签的 counter
var cc = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "url_interface_state_code",
	},
	[]string{"app"},
)

// Gauge metrics
func Gauge(domainList map[string]string) {
	// 注册指标
	EmptyRegistry.MustRegister(cc)

	for srvName, domainName := range domainList {

		res, err := http.Get("http://" + domainName)
		if err != nil {
			fmt.Println("[ERROR]: 访问接口异常,", err)
		} else {
			StatusCode := res.StatusCode
			if StatusCode == 200 {
				UrlStateCode = 1
			} else {
				UrlStateCode = 0
			}
		}

		// 写入标签
		cc.With(prometheus.Labels{
			"app": srvName,
		}).Set(float64(UrlStateCode))
	}

}

// RunServer 启动服务
func RunServer() {
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		promhttp.HandlerFor(EmptyRegistry,
			promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(writer, request)
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err)
	}
}

func main() {

	// 初始化配置文件
	config.Config()

	// 探测 Domain
	go Gauge(config.DomainMap)

	// 启动服务
	RunServer()

}
