package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/cinemaabyss/microservices/events/consumer"
	"github.com/cinemaabyss/microservices/events/model"
	"github.com/cinemaabyss/microservices/events/producer"
)

func main() {
	ctx := context.Background()

	port := getEnv("PORT", "8082")
	broker := getEnv("KAFKA_BROKERS", "127.0.0.1:9093")

	p := producer.NewProducer(broker)
	defer p.Close()

	handleHealth()
	handleMovieEvents(ctx, p)
	handleUserEvents(ctx, p)
	handlePaymentEvents(ctx, p)

	cons := consumer.NewConsumer(broker)
	cons.Run(ctx, "movie-events")
	cons.Run(ctx, "user-events")
	cons.Run(ctx, "payment-events")

	log.Printf("Starting events microservice on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func handleHealth() {
	http.HandleFunc("/api/events/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"status": true})
	})
}

func handleMovieEvents(ctx context.Context, p *producer.Producer) {
	http.HandleFunc("/api/events/movie", func(w http.ResponseWriter, r *http.Request) {
		var movieEvent model.MovieEvent

		if err := json.NewDecoder(r.Body).Decode(&movieEvent); err != nil {
			handleError(w, http.StatusBadRequest, err)
			return
		}

		e, err := p.CreateMovieEvent(ctx, movieEvent)

		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		handleSuccessCreated(w, e)
	})
}

func handleUserEvents(ctx context.Context, p *producer.Producer) {
	http.HandleFunc("/api/events/user", func(w http.ResponseWriter, r *http.Request) {
		var userEvent model.UserEvent

		if err := json.NewDecoder(r.Body).Decode(&userEvent); err != nil {
			handleError(w, http.StatusBadRequest, err)
			return
		}

		e, err := p.CreateUserEvent(ctx, userEvent)

		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		handleSuccessCreated(w, e)
	})
}

func handlePaymentEvents(ctx context.Context, p *producer.Producer) {
	http.HandleFunc("/api/events/payment", func(w http.ResponseWriter, r *http.Request) {
		var paymentEvent model.PaymentEvent

		if err := json.NewDecoder(r.Body).Decode(&paymentEvent); err != nil {
			handleError(w, http.StatusBadRequest, err)
			return
		}

		e, err := p.CreatePaymentEvent(ctx, paymentEvent)

		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		handleSuccessCreated(w, e)
	})
}

func handleError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	return
}

func handleSuccessCreated(w http.ResponseWriter, event model.Event) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(model.EventResponse{
		Status: "success",
		Event:  event,
	})
}
