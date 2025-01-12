package Kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

var (
	consumerKafkaBroker      = "localhost:9092"   // 替换为你的 Kafka 地址
	consumerKafkaTopic_Danmu = "xcxc_video_danmu" // 替换为你的 Kafka topic
	consumerGroupID          = "group1"           // 替换为你的消费者组 ID
)

func StartConsumer_Danmu() error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{consumerKafkaBroker},
		Topic:          consumerKafkaTopic_Danmu,
		GroupID:        consumerGroupID,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0,
	})

	defer reader.Close()

	log.Printf("Kafka consumer started, listening on topic %s\n", consumerKafkaTopic_Danmu)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		log.Printf("Received message: key=%s value=%s", string(msg.Key), string(msg.Value))
	}
}
