package reader

import (
	"fmt"
	"github.com/IBM/sarama"
)

type KafkaReaderBuilder struct {
	BrokersAddr []string
	Topic       string
}

func NewKafkaReaderBuilder(addr []string, topic string) *KafkaReaderBuilder {
	return &KafkaReaderBuilder{
		BrokersAddr: addr,
		Topic:       topic,
	}
}

func (k *KafkaReaderBuilder) Build() (Reader, error) {
	conf := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(k.BrokersAddr, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}
	return &KafkaReader{
		consumer: consumer,
		topic:    k.Topic,
	}, nil
}
