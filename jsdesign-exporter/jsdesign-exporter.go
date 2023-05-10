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
	"strings"
	"sync"
	"time"
)

var (
	// UrlStateCode url状态信息
	UrlStateCode = sync.Map{}
	// CertRemainingTime SSL 证书剩余时间
	CertRemainingTime = sync.Map{}
	// EmptyRegistry 清空默认指标
	EmptyRegistry = prometheus.NewRegistry()
	wg            sync.WaitGroup
	lock          sync.Mutex
)

// Monitor 指标采集器
type Monitor struct {
	InterfaceStatusCode  *prometheus.Desc
	SSLCertRemainingRime *prometheus.Desc
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

		SSLCertRemainingRime: prometheus.NewDesc(
			"ssl_cert_remaining_time",
			"Cert Remaining Time",
			[]string{"url"},
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
	descs <- m.SSLCertRemainingRime
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
		wg.Add(1)
		lock.Lock()
		go Gauge(srvName, domainName)
		lock.Unlock()
		wg.Wait()

		// 注册 url 状态指标
		srvNameValue, _ := UrlStateCode.Load(srvName)
		metrics <- prometheus.MustNewConstMetric(
			m.InterfaceStatusCode,
			prometheus.GaugeValue,
			float64(srvNameValue.(int)),
			srvName,
			domainName,
		)

		// 注册 SSL 证书有效期指标
		domainNameValue, _ := CertRemainingTime.Load(domainName)
		if domainNameValue != nil {
			protocolType := strings.Split(domainName, ":")
			if protocolType[0] == "https" {
				metrics <- prometheus.MustNewConstMetric(
					m.SSLCertRemainingRime,
					prometheus.GaugeValue,
					domainNameValue.(float64),
					domainName,
				)
			}
		}
	}

}

/*
	定义 HTTP client 配置
*/
var client = &http.Client{
	// 设置请求超时时间
	Timeout: 1 * time.Second,
	// 跳过安全检查
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

/*
	Gauge 方法
	用于获取状态信息, 并将信息状态返回给 map.
*/
func Gauge(srvName, domainName string) {

	defer wg.Done()

	resp, err := client.Get(domainName)
	if err != nil {
		global.GvaLogger.Sugar().Errorf("接口访问异常: %v", err.Error())
		UrlStateCode.Store(srvName, 0)
	} else {
		UrlStateCode.Store(srvName, 1)
	}

	if resp == nil {
		global.GvaLogger.Error("请求未响应")
		return
	}

	// 证书为空, 跳过检测
	if resp.TLS == nil {
		return
	}

	// 获取证书信息
	certs := resp.TLS.PeerCertificates[0]
	// 获取当前时间
	currentTime := time.Now().Unix()
	// 获取有效期时间
	certTime := certs.NotAfter.Unix()
	// 计算过期时间
	TimeRemaining := time.Unix(certTime, 0).Sub(time.Unix(currentTime, 0)).Seconds() / 86400
	CertRemainingTime.Store(domainName, TimeRemaining)
	//fmt.Println("证书有效期：", certs.NotBefore.Format(time.RFC3339), "-", certs.NotAfter.Format(time.RFC3339))

}

/*
	RunServer 服务启动
*/
func RunServer() {

	global.GvaLogger.Info("Server Started Successful.")

	// 注册指标
	//fmt.Println(NewMonitorMetrics())
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
