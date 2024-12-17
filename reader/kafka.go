package reader

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

type KafkaReader struct {
	consumer sarama.Consumer
	topic    string
}

func (k *KafkaReader) Read(ctx context.Context, ch chan<- []byte) error {
	// 获取所有分区
	partitionList, err := k.consumer.Partitions(k.topic)
	if err != nil {
		return fmt.Errorf("get partitions error: %v", err)
	}

	var wg sync.WaitGroup
	for _, partition := range partitionList {
		// 针对每个分区创建一个对应的分区消费者
		pc, err := k.consumer.ConsumePartition(k.topic, partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d: %v\n", partition, err)
			continue // 继续尝试消费其他分区
		}

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			defer pc.AsyncClose()

			// 消费消息
			for {
				select {
				case msg := <-pc.Messages():
					ch <- msg.Value
				case <-ctx.Done():
					// 收到停止信号，退出 goroutine
					return
				}
			}
		}(pc)
	}
	// 等待所有 goroutine 完成
	wg.Wait()
	return nil
}
