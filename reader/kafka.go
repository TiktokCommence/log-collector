package reader

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

// KafkaReader 结构体
type KafkaReader struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	groupID       string
}

// newKafkaReader 初始化 KafkaReader
func newKafkaReader(brokers []string, topic, groupID string) (*KafkaReader, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true // 启用自动提交偏移量

	// 创建消费者组
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %v", err)
	}

	return &KafkaReader{
		consumerGroup: consumerGroup,
		topic:         topic,
		groupID:       groupID,
	}, nil
}

// Kafka 消费者处理器
type messageHandler struct {
	ch chan<- []byte
}

// Setup 初始化消费者
func (h *messageHandler) Setup(sess sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 清理资源
func (h *messageHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费消息
func (h *messageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 处理消息
		h.ch <- msg.Value
		// 提交偏移量，表示消息已被消费
		sess.MarkMessage(msg, "")
	}
	return nil
}

// Read 从 Kafka 中读取消息
func (k *KafkaReader) Read(ctx context.Context, ch chan<- []byte) error {
	handler := &messageHandler{ch: ch}

	// 启动消费者组
	for {
		// 这里传递的参数包括消费者组的上下文（以便控制退出）、消费者组 ID 和消费的 topic
		err := k.consumerGroup.Consume(ctx, []string{k.topic}, handler)
		if err != nil {
			// 发生错误时打印日志并返回
			log.Printf("Error consuming messages: %v", err)
		}

		// 检查是否已收到停止信号（退出条件）
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
