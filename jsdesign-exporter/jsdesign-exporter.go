package main

import (
	"crypto/tls"
	"exporter/config"
	"exporter/global"
	"exporter/initialize"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	// UrlStateCode url状态码
	UrlStateCode int
	// EmptyRegistry 清空默认指标
	EmptyRegistry = prometheus.NewRegistry()
	// 带动态标签的 counter
	cc = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "url_interface_state_code",
		},
		[]string{"app", "url"},
	)
)

// Gauge metrics
func Gauge(srvName, domainName string) {

	client := &http.Client{
		// 设置请求超时时间
		Timeout: 1 * time.Second,
		// 跳过安全检查
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	_, err := client.Get(domainName)
	if err != nil {
		global.GvaLogger.Sugar().Errorf("接口访问异常: %v", err.Error())
		UrlStateCode = 0
	} else {
		UrlStateCode = 1
	}

	// 写入标签
	cc.With(prometheus.Labels{
		"app": srvName,
		"url": domainName,
	}).Set(float64(UrlStateCode))

}

func RunServer() {

	global.GvaLogger.Info("Server Started Successful.")

	// 注册指标
	EmptyRegistry.MustRegister(cc)

	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {

		global.GvaLogger.Info("--- 开始抓取指标! ---")
		for srvName, domainName := range config.DomainMap {

			// 探测 Domain 状态
			go Gauge(srvName, domainName)

		}

		promhttp.HandlerFor(EmptyRegistry,
			promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(writer, request)
	})

	if err := http.ListenAndServe(":"+global.GvaServerConfig.Exporter.Port, nil); err != nil {
		fmt.Println(err)
	}

}

func main() {

	// 初始化配置文件
	initialize.InitConfig()
	config.Config()
	// 初始化日志
	initialize.InitLogger()

	// 启动服务
	RunServer()

}
