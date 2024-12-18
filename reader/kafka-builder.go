package reader

import "fmt"

const GROUPID = "appLog"

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
	r, err := newKafkaReader(k.BrokersAddr, k.Topic, GROUPID)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka %w", err)
	}
	return r, nil
}
