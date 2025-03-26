package types

// KafkaLogField kafka 日志结构
type KafkaLogField struct {
	Timestamp string `json:"@timestamp"`
	Namespace string `json:"k8s_pod_namespace"`
	Service   string `json:"k8s_container_name"`
	Message   string `json:"message"`
}

type Message struct {
	Level string `json:"level"`
}

func (m Message) GetLevel() string {
	if m.Level == "" {
		m.Level = "Unknown"
	}
	return m.Level
}

// ServiceId 服务唯一 ID
type ServiceId uint64

// MetricLabel metric 标签
type MetricLabel struct {
	// 命名空间
	Namespace string
	// 服务
	Service string
	// 日志等级
	Level string
	// 值
	Value float64
}
