package domain

import "time"

type ContactRequest struct {
	Name    string `json:"name" validate:"required,min=2"`
	Email   string `json:"email" validate:"required,email"`
	Phone   string `json:"phone" validate:"required"`
	Comment string `json:"comment" validate:"required,min=5"`
}

type AIResult struct {
	Sentiment string `json:"sentiment"`
	Reply     string `json:"reply"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	IP        string `json:"ip"`
	Status    int    `json:"status"`
}

type SentimentStats struct {
	Positive int `json:"positive" example:"8"`
	Neutral  int `json:"neutral" example:"3"`
	Negative int `json:"negative" example:"1"`
}

type Metrics struct {
	TotalRequests  int            `json:"total_requests" example:"12"`
	SentimentStats SentimentStats `json:"sentiment_stats"`
	LastUpdate     time.Time      `json:"last_update" example:"2026-07-17T13:40:00Z"`
}

type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

type ContactResponse struct {
	Success   bool   `json:"success" example:"true"`
	Sentiment string `json:"sentiment" example:"positive"`
	AiReply   string `json:"ai_reply" example:"Здравствуйте! Спасибо за ваше обращение..."`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Описание ошибки бэкенда"`
}
