package writer

import "github.com/segmentio/kafka-go"

func CreateNewWriter(broker, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}
}
