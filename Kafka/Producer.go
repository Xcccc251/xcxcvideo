package Kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

// Kafka
const (
	kafkaBroker = "localhost:9092"
)

// 生产者逻辑
func ProduceMessage(topic string, message []byte) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer writer.Close()

	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", time.Now().UnixNano())),
			Value: message,
		},
	)

	return err
}
