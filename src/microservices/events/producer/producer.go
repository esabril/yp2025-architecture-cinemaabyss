package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cinemaabyss/microservices/events/model"
	"github.com/cinemaabyss/microservices/events/writer"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	moviesWriter   *kafka.Writer
	usersWriter    *kafka.Writer
	paymentsWriter *kafka.Writer
}

func NewProducer(broker string) *Producer {
	return &Producer{
		moviesWriter:   writer.CreateNewWriter(broker, "movie-events"),
		usersWriter:    writer.CreateNewWriter(broker, "user-events"),
		paymentsWriter: writer.CreateNewWriter(broker, "payment-events"),
	}
}

func (p *Producer) CreateMovieEvent(ctx context.Context, data model.MovieEvent) (model.Event, error) {
	e := model.Event{
		ID:        "new-movie-event-" + time.Now().Format(time.RFC3339Nano),
		Type:      "movie",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   data,
	}

	return e, p.publish(ctx, p.moviesWriter, e)
}

func (p *Producer) CreateUserEvent(ctx context.Context, data model.UserEvent) (model.Event, error) {
	e := model.Event{
		ID:        "new-user-event-" + time.Now().Format(time.RFC3339Nano),
		Type:      "user",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   data,
	}

	return e, p.publish(ctx, p.usersWriter, e)
}

func (p *Producer) CreatePaymentEvent(ctx context.Context, data model.PaymentEvent) (model.Event, error) {
	e := model.Event{
		ID:        "new-payment-event-" + time.Now().Format(time.RFC3339Nano),
		Type:      "payment",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload:   data,
	}

	return e, p.publish(ctx, p.paymentsWriter, e)
}

func (p *Producer) Close() {
	p.moviesWriter.Close()
	p.usersWriter.Close()
	p.paymentsWriter.Close()
}

func (p *Producer) publish(ctx context.Context, writer *kafka.Writer, event model.Event) error {
	e, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if err := writer.WriteMessages(ctx, kafka.Message{Value: e}); err != nil {
		return err
	}

	return nil
}
