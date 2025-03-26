package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/logc"
	"log2metric/client"
	"log2metric/collector"
	"log2metric/storage"
	"log2metric/types"
	"net/http"
	"os"
	"strings"
	"time"
)

type Env struct {
	TopicPrefix string
	Topics      []string
	Brokers     []string
}

func GetEnvInfo() (*Env, error) {
	topicPrefix := os.Getenv("TOPIC_PREFIX")
	topics := os.Getenv("TOPICS")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	if kafkaBrokers == "" {
		return nil, fmt.Errorf("请配置 KAFKA_BROKERS, 例如: 1.1.1.1:9092,2.2.2.2:9092...")
	}

	var topicList []string
	if topics == "" && topicPrefix == "" {
		return nil, fmt.Errorf("Topics 和 TopicPrefix 不能同时为空, 请指定具体的 Topics(TOPICS) 或 Topic 前缀(TOPIC_PREFIX)")
	}

	if topics != "" {
		if !strings.Contains(topics, ",") {
			topicList = []string{topics}
		} else {
			topicList = strings.Split(topics, ",")
		}
	}

	var brokerList []string
	if !strings.Contains(kafkaBrokers, ",") {
		brokerList = []string{kafkaBrokers}
	} else {
		brokerList = strings.Split(kafkaBrokers, ",")
	}

	return &Env{
		TopicPrefix: topicPrefix,
		Topics:      topicList,
		Brokers:     brokerList,
	}, nil
}

func (e Env) SetBrokers(config *client.KafkaConfig) {
	config.Brokers = e.Brokers
}

func (e Env) GetTopics(config *client.KafkaConfig) error {
	if e.Topics != nil {
		config.Topics = e.Topics
		return nil
	}

	err := client.GetTopicsByPrefix(e.TopicPrefix, config)
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	data := storage.NewCacheStorage()
	config := &client.KafkaConfig{
		Config: sarama.NewConfig(),
		//Brokers:        []string{"172.17.84.75:9092", "172.17.84.76:9092", "172.17.84.77:9092"},
		GroupID:        "log2metric-group",
		CommitInterval: 1 * time.Second,
	}

	info, err := GetEnvInfo()
	if err != nil {
		panic(err)
	}
	info.SetBrokers(config)
	err = info.GetTopics(config)
	if err != nil {
		panic(err)
	}

	go kafkaConsumer(data, config)

	metric := collector.NewProcessGaugeCollect(data)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metric)

	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "hello")
	})

	http.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		promhttp.HandlerFor(registry,
			promhttp.HandlerOpts{ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(writer, request)
	})

	logc.Info(context.Background(), "Service started!")

	if err := http.ListenAndServe(":"+"9099", nil); err != nil {
		panic(err)
	}
}

func generateHash(namespace, service, level string) uint64 {
	return hash.Hash([]byte(fmt.Sprintf("%s-%s-%s", namespace, service, level)))
}

func kafkaConsumer(data *storage.CacheStorage, config *client.KafkaConfig) {
	consumer, err := client.NewKafkaConsumer(config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	handler := func(msg *sarama.ConsumerMessage) {
		var logField types.KafkaLogField
		err := json.Unmarshal(msg.Value, &logField)
		if err != nil {
			logc.Error(context.Background(), err.Error())
			return
		}

		var messageField types.Message
		err = json.Unmarshal([]byte(logField.Message), &messageField)
		if err != nil {
			logc.Alert(context.Background(), fmt.Sprintf("Message 格式不合法, service: %s, err: %s, msg: %s", logField.Service, err.Error(), logField.Message))
			return
		}

		level := strings.ToUpper(messageField.GetLevel())
		// 计算服务唯一 ID
		hashId := generateHash(logField.Namespace, logField.Service, level)
		data.Set(types.ServiceId(hashId), types.MetricLabel{
			Namespace: logField.Namespace,
			Service:   logField.Service,
			Level:     level,
		})
	}

	if err := consumer.Consume(handler); err != nil {
		panic(err)
	}

	select {}
}
