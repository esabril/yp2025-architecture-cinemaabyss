package model

import "time"

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp string      `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type EventResponse struct {
	Status    string `json:"status"`
	Partition int    `json:"partition"`
	Offset    int    `json:"offset"`
	Event     Event  `json:"event"`
}

type MovieEvent struct {
	Id          int64    `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"`
	UserId      int64    `json:"user_id"`
	Rating      string   `json:"rating"`
	Genres      []string `json:"genres"`
	Description string   `json:"description"`
}

type UserEvent struct {
	UserId    int       `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentId  int       `json:"payment_id"`
	UserId     int       `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	MethodType string    `json:"method_type"`
}
