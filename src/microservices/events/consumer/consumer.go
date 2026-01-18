package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	brokers []string
}

func NewConsumer(broker string) *Consumer {
	return &Consumer{
		brokers: []string{broker},
	}
}

func (c Consumer) Run(ctx context.Context, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		Topic:   topic,
		GroupID: "events-service",
	})

	log.Printf("📥 Запущен consumer (%s)", topic)

	go func() {
		defer func() { _ = r.Close() }()
		for {
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				log.Printf("Topic %s consumer error: %v", topic, err)
				continue
			}

			log.Printf("📩 Получено событие:\ntopic: %s, event: %s\n", msg.Topic, msg.Value)
		}
	}()
}
