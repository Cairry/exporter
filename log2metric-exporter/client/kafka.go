package client

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/zeromicro/go-zero/core/logc"
	"strings"
	"sync/atomic"
	"time"
)

import (
	"context"
)

type KafkaConfig struct {
	Config         *sarama.Config
	Brokers        []string // Kafka集群地址
	Topics         []string
	GroupID        string
	CommitInterval time.Duration // 手动提交间隔（0表示禁用）
}

type TimestampConsumer struct {
	group         sarama.ConsumerGroup
	client        sarama.Client
	config        *KafkaConfig
	lastCommit    atomic.Value
	consumerCtx   context.Context
	cancelConsume context.CancelFunc
}

// GetTopicsByPrefix 通过前缀获取 Topics
func GetTopicsByPrefix(prefix string, config *KafkaConfig) error {
	admin, err := sarama.NewClusterAdmin(config.Brokers, config.Config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create admin client: %v", err))
	}
	defer admin.Close()

	// 获取所有 Topic 元数据
	listTopicsResp, err := admin.ListTopics()
	if err != nil {
		return fmt.Errorf("list topics failed: %w", err)
	}

	// 过滤匹配前缀的 Topic
	var matched []string
	for topic := range listTopicsResp {
		if strings.HasPrefix(topic, prefix) {
			matched = append(matched, topic)
		}
	}

	config.Topics = matched

	return nil
}

func NewKafkaConsumer(config *KafkaConfig) (*TimestampConsumer, error) {
	// 配置Sarama客户端
	cfg := config.Config
	cfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	cfg.Consumer.Offsets.AutoCommit.Enable = false // 强制禁用自动提交
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 创建Sarama客户端
	client, err := sarama.NewClient(config.Brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	// 创建消费者组
	group, err := sarama.NewConsumerGroupFromClient(config.GroupID, client)
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &TimestampConsumer{
		group:  group,
		client: client,
		config: config,
	}, nil
}

// Consume 开始消费指定的主题
func (c *TimestampConsumer) Consume(handler func(*sarama.ConsumerMessage)) error {
	c.consumerCtx, c.cancelConsume = context.WithCancel(context.Background())

	handlerWrapper := &timestampHandler{
		client:      c.client,
		msgHandler:  handler,
		commitTimer: time.NewTicker(c.config.CommitInterval),
	}

	go func() {
		for {
			// 监听上下文取消信号
			select {
			case <-c.consumerCtx.Done():
				return
			default:
				err := c.group.Consume(c.consumerCtx, c.config.Topics, handlerWrapper)
				if err != nil {
					logc.Error(c.consumerCtx, fmt.Sprintf("Consumer failed, err: %s", err.Error()))
					continue
				}
			}
		}
	}()

	return nil
}

// Close 关闭消费者，释放资源
func (c *TimestampConsumer) Close() error {
	if c.cancelConsume != nil {
		c.cancelConsume()
	}
	if err := c.group.Close(); err != nil {
		return err
	}
	return c.client.Close()
}

// timestampHandler 处理Kafka消息的处理器
type timestampHandler struct {
	client      sarama.Client
	targetTime  int64
	msgHandler  func(*sarama.ConsumerMessage)
	commitTimer *time.Ticker
}

// Setup 在消费者启动时调用，重置偏移量到最新
func (h *timestampHandler) Setup(session sarama.ConsumerGroupSession) error {
	for topic, partitions := range session.Claims() {
		for _, partition := range partitions {
			offset, err := h.client.GetOffset(topic, partition, sarama.OffsetNewest)
			if err != nil {
				return err
			}
			session.MarkOffset(topic, partition, offset, "")
			logc.Info(context.Background(), fmt.Sprintf("Reset %s-%d to offset %d", topic, partition, offset))
		}
	}
	return nil
}

// Cleanup 在消费者关闭时调用，停止定时器
func (h *timestampHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.commitTimer.Stop()
	return nil
}

// ConsumeClaim 处理消息并提交偏移量
func (h *timestampHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			h.msgHandler(msg)
			session.MarkMessage(msg, "")
		case <-h.commitTimer.C:
			session.Commit()
		}
	}
}
