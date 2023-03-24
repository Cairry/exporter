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
	// UrlStateCode url状态信息
	UrlStateCode = make(map[string]int)
	// EmptyRegistry 清空默认指标
	EmptyRegistry = prometheus.NewRegistry()
)

// Monitor 指标采集器
type Monitor struct {
	InterfaceStatusCode *prometheus.Desc
}

/*
	NewMonitorMetrics 指标采集器规范
	定义指标相关信息:
		fqName: 指标名称
		help: 帮助描述信息
		variableLabels: 动态label名称数组
		constLabels: labels
*/
func NewMonitorMetrics() *Monitor {
	return &Monitor{
		InterfaceStatusCode: prometheus.NewDesc(
			"url_interface_state_code",
			"url_interface_state_code",
			[]string{"app", "url"},
			nil,
		),
	}
}

/*
	Describe 方法 收集描述信息
	实现 Collector 接口
	用于传递所有可能的指标描述信息, 可以在程序运行期间添加新的描述, 收集信息的指标信息.
*/
func (m Monitor) Describe(descs chan<- *prometheus.Desc) {
	//TODO implement me
	descs <- m.InterfaceStatusCode
}

/*
	Collect 方法 实施指标抓取
	实现 Collector 接口
	用于注册器调用Collect执行实际的抓取指标工作, 并将收集的数据传递到Channel中返回.
	收集的指标信息来自雨Describe方法中传递, 可以并发执行抓取工作, 但是必须保证线程的安全.
*/
func (m Monitor) Collect(metrics chan<- prometheus.Metric) {
	//TODO implement me
	for srvName, domainName := range config.DomainMap {
		// 探测 Domain 状态
		go Gauge(srvName, domainName)

		metrics <- prometheus.MustNewConstMetric(
			m.InterfaceStatusCode,
			prometheus.GaugeValue,
			float64(UrlStateCode[srvName]),
			srvName,
			domainName,
		)
	}
}

/*
	Gauge 方法
	用于获取状态信息, 并将信息状态返回给 map.
*/
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
		UrlStateCode[srvName] = 0
	} else {
		UrlStateCode[srvName] = 1
	}
}

/*
	RunServer 服务启动
*/
func RunServer() {

	global.GvaLogger.Info("Server Started Successful.")

	// 注册指标
	fmt.Println(NewMonitorMetrics())
	EmptyRegistry.MustRegister(NewMonitorMetrics())

	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		global.GvaLogger.Info("--- 开始抓取指标! ---")

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
